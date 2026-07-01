# Email over XNS

Email can run on top of XNS when the name owner publishes mail records through owner DNS and exposes the usual mail ports on the same I2P destination.

The route is ordinary at the application layer:

```text
user@example.xns
    |
MX lookup for example.xns
    |
mail.example.xns
    |
XNS Resolver synthetic address
    |
I2P destination derived from the owner key
    |
mail server
```

The mail server does not need to understand XNS. It needs to listen on the ports that i2pd forwards, and clients need to resolve `.xns` names through XNS Resolver.

## Requirements

This guide assumes [DNS over XNS](/docs/guides/dns-over-xns) is already working.

You will need:

- an owner DNS server on port `53`
- i2pd server tunnels for the mail ports
- a mail server
- XNS Resolver on the client side

The examples use `example.xns`, `mail.example.xns` and maddy.

## I2P tunnels

Add mail ports to the same i2pd destination used by the XNS name:

```ini
[example_smtp]
type = server
host = 127.0.0.1
inport = 25
port = 9025
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5

[example_submission]
type = server
host = 127.0.0.1
inport = 587
port = 9587
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5

[example_submission_tls]
type = server
host = 127.0.0.1
inport = 465
port = 9465
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5

[example_imap]
type = server
host = 127.0.0.1
inport = 993
port = 9993
keys = xns/private.dat
signaturetype = 11
i2cp.leaseSetType = 5
```

Port `25` is for server-to-server SMTP. Port `587` is for authenticated submission with STARTTLS. Port `465` is implicit TLS submission. Port `993` is IMAP over TLS.

Restart i2pd after changing the tunnel configuration:

```sh
sudo systemctl restart i2pd
```

## DNS records

Add mail records to the owner DNS zone:

```dns
@       IN MX   10 mail.example.xns.
@       IN TXT  "v=spf1 -all"

_dmarc  IN TXT  "v=DMARC1; p=reject; adkim=s; aspf=s"
```

The MX record tells mail software to deliver mail for `example.xns` to `mail.example.xns`.

The SPF record is deliberately strict. IP-based authorization is not the trust boundary here. The useful proof is DKIM: outgoing mail is signed by a key published under the XNS name, and receivers verify that signature through owner DNS. DMARC then tells receivers to reject mail that does not align with that domain policy.

## maddy

maddy is a small all-in-one mail server. It can provide SMTP, submission, IMAP, local accounts and DKIM signing from one config.

Create a local state directory and a TLS certificate:

```sh
mkdir -p /srv/xns-mail/data /srv/xns-mail/run /srv/xns-mail/tls

openssl req -x509 -newkey rsa:4096 -nodes -days 3650 \
  -subj "/CN=mail.example.xns" \
  -keyout /srv/xns-mail/tls/privkey.pem \
  -out /srv/xns-mail/tls/fullchain.pem
```

The certificate may be self-signed. I2P already provides encrypted transport to the destination, but mail clients still expect TLS on submission and IMAP. A self-signed certificate will have to be accepted by the client.

A local XNS setup should bind only to loopback. The public side is the I2P tunnel:

```maddy
$(hostname) = mail.example.xns
$(primary_domain) = example.xns
$(local_domains) = $(primary_domain)
$(base_dir) = /srv/xns-mail

state_dir $(base_dir)/data
runtime_dir $(base_dir)/run

tls file $(base_dir)/tls/fullchain.pem $(base_dir)/tls/privkey.pem

auth.pass_table local_authdb {
    table sql_table {
        driver sqlite3
        dsn $(base_dir)/data/credentials.db
        table_name passwords
    }
}

storage.imapsql local_mailboxes {
    driver sqlite3
    dsn $(base_dir)/data/imapsql.db
}

hostname $(hostname)

msgpipeline local_routing {
    destination postmaster $(local_domains) {
        deliver_to &local_mailboxes
    }

    default_destination {
        reject 550 5.1.1 "User doesn't exist"
    }
}

smtp tcp://127.0.0.1:9025 {
    dmarc yes
    check {
        dkim
        spf
    }

    source $(local_domains) {
        reject 501 5.1.8 "Use Submission for outgoing SMTP"
    }

    default_source {
        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }
        default_destination {
            reject 550 5.1.1 "User doesn't exist"
        }
    }
}

submission tls://127.0.0.1:9465 tcp://127.0.0.1:9587 {
    auth &local_authdb

    source $(local_domains) {
        check {
            authorize_sender {
                user_to_email identity
            }
        }

        destination postmaster $(local_domains) {
            deliver_to &local_routing
        }

        default_destination {
            modify {
                dkim $(primary_domain) $(local_domains) default
            }
            deliver_to &remote_queue
        }
    }

    default_source {
        reject 501 5.1.8 "Non-local sender domain"
    }
}

target.remote outbound_delivery {
    mx_auth {
        local_policy {
            min_tls_level none
            min_mx_level none
        }
    }
}

target.queue remote_queue {
    target &outbound_delivery
    autogenerated_msg_domain $(primary_domain)
}

imap tls://127.0.0.1:9993 {
    auth &local_authdb
    storage &local_mailboxes
}
```

Verify the config:

```sh
maddy --config maddy.conf verify-config
```

This also creates the DKIM key used by the `dkim $(primary_domain) $(local_domains) default` modifier. With the state directory above, maddy writes:

```text
/srv/xns-mail/data/dkim_keys/example.xns_default.key
/srv/xns-mail/data/dkim_keys/example.xns_default.dns
```

The `.key` file is the private signing key. The `.dns` file contains the public TXT value that must be published in owner DNS:

```sh
cat /srv/xns-mail/data/dkim_keys/example.xns_default.dns
```

Add that value to the zone under `default._domainkey`:

```dns
default._domainkey IN TXT "v=DKIM1; k=rsa; p=<public-key>"
```

Long DKIM keys may need to be split into multiple quoted strings:

```dns
default._domainkey IN TXT (
    "v=DKIM1; k=rsa; p=first-part"
    "second-part"
)
```

DNS joins those strings into one TXT value. Restart or reload the DNS server after changing the zone.

Create an account:

```sh
maddy --config maddy.conf creds create postmaster@example.xns
maddy --config maddy.conf imap-acct create postmaster@example.xns
```

Start the server:

```sh
maddy --config maddy.conf run
```

## Testing

Check the DNS records through XNS Resolver:

```sh
dig example.xns MX +short
dig default._domainkey.example.xns TXT +short
dig _dmarc.example.xns TXT +short
```

Check that the mail host resolves through the virtual route:

```sh
getent hosts mail.example.xns
```

Then check the mail ports:

```sh
openssl s_client \
  -connect mail.example.xns:587 \
  -starttls smtp \
  -servername mail.example.xns

openssl s_client \
  -connect mail.example.xns:993 \
  -servername mail.example.xns
```

The certificate may be self-signed. That warning is expected unless the operator uses a certificate trusted by the client.

## Mail clients

Automatic account discovery may not understand `.xns`. Configure the account manually:

```text
Email address: user@example.xns
Username:      user@example.xns

Incoming IMAP:
Host:     mail.example.xns
Port:     993
Security: SSL/TLS

Outgoing SMTP:
Host:     mail.example.xns
Port:     587
Security: STARTTLS
```

The client must use the system resolver path where XNS Resolver is active. If it bypasses the system resolver or runs inside a sandbox that cannot see the XNS route, it will not find `mail.example.xns`.

## What is proven

XNS establishes which owner key controls `example`. I2P reaches the destination derived from that key. Owner DNS publishes the MX, DKIM and DMARC records. DKIM signs the outgoing mail for the domain.

This is not ordinary global clearnet mail discovery. It is mail inside an XNS-aware environment, using ordinary mail protocols after resolution has taken the XNS route.
