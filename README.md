# Hummingbard

Hummingbard is an experimental client for building decentralized communities on top of [Matrix](https:/matrix.org). See a live instance on [hummingbard.com](https://hummingbard.com)


### What Works
- Register local and federated users
- Federated logins with existing Matrix accounts
- Join local and federated spaces
- Follow local and federated users
- Generic post editor (markdown)
- Quick posts, with images/attachments/links/youtube/etc
- Blog posts with slug/metadata
- Replies to posts
- Sharing posts on profiles and across spaces
- User feed
- Public feed
- Create local and federated spaces 
- Different space types - community, gallery
- Customize spaces and user profiles with basic info, custom CSS
- Deeply nested spaces (`/music/jazz/fusion`)

### What Doesn't Work
- Private spaces and user profiles
- Embedded chat in spaces
- Direct Messages
- Registration flows


### Dendrite
Hummingbard relies on these features that are currently only implemented in Dendrite, or expected to be implemented soon:

- Spaces ([MSC2946](https://github.com/matrix-org/matrix-doc/pull/2946))
- Threading ([MSC2836](https://github.com/matrix-org/matrix-doc/pull/2836))

There is a temporary patch in our [forked
Dendrite](https://github.com/hummingbard/dendrite) for paginating threads. This
should not be necessary once upstream Dendrite implements threads fully.

## Install

To run Hummingbard, you'll need:

- [Dendrite fork](https://github.com/hummingbard/dendrite) configured and running
- redis (for session storage)
- postgres (for various non-Matrix storage)
- [goose](https://github.com/pressly/goose) for migrations
- uglifycss - `npm install uglifycss`

### Steps:

1. Clone the repo
2. Copy `config-sample.toml` to `config.toml`, update with DB config etc.
3. Run `make` 
4. Run migrations in `db/migrations`
5. Run `npm run build` in `/ui/js`
6. Pull a JSON dump for large matrix rooms with `curl 'https://matrix-client.matrix.org:443/_matrix/client/r0/publicRooms?limit=500' > bigrooms.json` (we avoid large rooms to help Dendrite not consume too much resources)
7. Run the binary `./bin/hummingbard`

### You may want to:
1. Put Hummingbard behind Nginx
2. Server static files via Nginx
3. Use a systemd unit if appropriate

## License
The code is currenly licensed under [AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html). I may choose a more permissive license in the future.
