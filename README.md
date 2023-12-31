## fonk

This is a personal fork of honk, nothing here is
stable and anything can break at any time!  
You probably want to visit the original project instead:
https://humungus.tedunangst.com/r/honk

### Features

Take control of your honks and join the federation.  
An ActivityPub server with minimal setup and support costs.  
Spend more time using the software and less time operating it.

No attention mining.  
No likes, no faves, no polls, no stars, no claps, no counts.

Purple color scheme. Custom emus. Memes too.  
Avatars automatically assigned by the NSA.

The button to submit a new honk says "it's gonna be honked".

The honk mission is to work well if it's what you want.  
This does not imply the goal is to be what you want.

### Guidelines

One honk per day, or call it an "eighth-tenth" honk.  
If your honk frequency changes, so will the number of honks.

The honk should be short, but not so short that you cannot identify it.

The honk is an animal sign of respect and should be accompanied by a  
friendly greeting or a nod.

The honk should be done from a seat and in a safe area.

It is considered rude to make noise in a place of business.

The honk may be made on public property only when the person doing  
the honk has the permission of the owner of that property.

### Build

To build honk, you will need a go compiler version 1.21 or later.

```sh
git clone https://github.com/nothub/fonk.git
cd honk && make
```

Even on a fast machine, building from source can take several seconds.

### Setup

honk expects to be fronted by a TLS terminating reverse proxy.

First, create the database:

```sh
./honk init [--hash <hash>] [--listen <host:port>] <username> <fqdn>
```

If no hash (bcrypt) is supplied, an interactive password dialog awaits your input.

The listen address can be a tcp or unix socket, e.g. `127.0.0.1:31337` or `/var/www/honk.sock`.  
The listen address defaults to: `0.0.0.0:8080`

Then run honk: `./honk`

### Upgrade

Prior to any honk upgrade, you should use the `backup` command:

```sh
./old-honk backup "$(date +backup-%F)"
```

Then with the new honk release, just run the `upgrade` command:

```sh
./new-honk upgrade
./new-honk
```

### Docker

honk is available packaged as a
[Docker image](https://hub.docker.com/r/n0thub/fonk).

<details>
  <summary>Usage examples</summary>

##### Persistent data volume

```sh
docker run --rm            \
  -p "127.0.0.1:8080:8080" \
  -v "${PWD}/data:/data"   \
  "n0thub/fonk:latest"
```

---

##### Database setup

The database will be initialized if not found.  
A password can be supplied interactive or by flag as bcrypt hash.

Handy ways to generate a bcrypt hash:
- `hash="$(./honk genhash)"`
- `hash="$(htpasswd -nBC 12 "" | tr -d ':\n')"`

```sh
hash="$(./honk genhash)"
docker run --rm          \
  -v "${PWD}/data:/data" \
  "n0thub/fonk:latest"   \
    init                 \
    --hash "${hash}"     \
    "admin"              \
    "honk.example.org"
```

---

##### Database upgrade

A database upgrade can be executed by passing the required command to the
container.

```sh
docker run --rm              \
  -v "${PWD}/data:/data"     \
  "n0thub/fonk:latest"       \
  "upgrade"
```

---

##### Custom HTML views

```sh
docker run --rm               \
  -p "127.0.0.1:8080:8080"    \
  -v "${PWD}/data:/data"      \
  -v "${PWD}/views:/views:ro" \
  "n0thub/fonk:latest"
```

---

##### Custom UID & GID

```sh
docker run --rm            \
  -p "127.0.0.1:8080:8080" \
  -v "${PWD}/data:/data"   \
  -e "PUID=9001"           \
  -e "PGID=9002"           \
  "n0thub/fonk:latest"
```

</details>

### Documentation

There is a [more complete incomplete manual](./docs/). This is just the README.

### Disclaimer

Do not use honk to contact emergency services.  
