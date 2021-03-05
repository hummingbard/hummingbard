let tl = document.querySelector(`.timeline`)
if(tl) {
  import('./timeline.svelte').then(res => {
      new res.default({
          target: tl,
          props: {
          },
          hydrate: true
      });
  })
}

let np = document.querySelector(`.new-post`)
if(np && authenticated) {
  import('../new-post/new-post.svelte').then(res => {
      new res.default({
          target: np,
          props: {
          },
          hydrate: true
      });
  })
}

let replies = document.querySelector(`.replies`)
if(replies) {
  import('./replies.svelte').then(res => {
      new res.default({
          target: replies,
          props: {
          },
          hydrate: true
      });
  })
}

let ngp = document.querySelector(`.new-gallery-post`)
if(ngp) {
  import('../gallery/new-gallery-post.svelte').then(res => {
      new res.default({
          target: ngp,
          props: {
          },
          hydrate: true
      });
  })
}

let gal = document.querySelector(`.gallery`)
if(gal) {
  import('../gallery/gallery.svelte').then(res => {
      new res.default({
          target: gal,
          props: {
          },
          hydrate: true
      });
  })
}

let crp = document.querySelector(`.create-post`)
if(crp) {
  import('../new-post/view/view.svelte').then(res => {
      new res.default({
          target: crp,
          props: {
          },
          hydrate: true
      });
  })
}

let edpg = document.querySelector(`.edit-page`)
if(edpg) {
  import('../page/edit-page.svelte').then(res => {
      new res.default({
          target: edpg,
          props: {
          },
          hydrate: true
      });
  })
}

let lm = document.querySelector(`.load-more`)
if(lm) {
  if(window.timeline.end || window.timeline.userFeed) {
    loadMore()
  }
}

let sorr = document.querySelector(`.sort-replies`)
if(sorr) {
  import('../timeline/sort-replies.svelte').then(res => {
      new res.default({
          target: sorr,
          props: {
          },
          hydrate: true
      });
  })
}


let More;
function loadMore() {
  import('./more.svelte').then(res => {
    More = res.default
      new More({
          target: lm,
          props: {
            type: window.timeline.type,
          },
          hydrate: true
      });
  })
}

let ap = document.querySelector(`.added-posts`)
if(ap) {
  import('./added-posts.svelte').then(res => {
      new res.default({
          target: ap,
          props: {
          },
          hydrate: true
      });
  })
}

let joins = document.querySelectorAll(`.join-room`)
if(joins) {
  import('../join/join.svelte').then(res => {
    joins.forEach(join => {
      new res.default({
          target: join,
          props: {
            type: join.dataset.type,
            alias: join.dataset.alias,
            name: join.dataset.name,
            id: join.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

let youtubes = document.querySelectorAll(`.youtube`)
if(youtubes) {
  import('../post/link/youtube.svelte').then(res => {
    youtubes.forEach(youtube => {
      new res.default({
          target: youtube,
          props: {
            id: youtube.dataset.id,
            title: youtube.dataset.title,
            description: youtube.dataset.description,
            href: youtube.dataset.href,
          },
          hydrate: true
      });
    })
  })
}

let ts = document.querySelector(`.room-settings`)
if(ts && window?.timeline?.admin) {
  import('./settings/settings.svelte').then(res => {
      new res.default({
          target: ts,
          hydrate: true
      });
  })
}

let rmem = document.querySelector(`.room-members`)
if(rmem && authenticated && window?.timeline?.member) {
  import('../timeline/members/members.svelte').then(res => {
      new res.default({
          target: rmem,
          props: {
            members: rmem.dataset.members,
          },
          hydrate: true
      });
  })
}

let lrs = document.querySelector(`.load-room-settings`)
if(lrs && window?.timeline?.admin) {
  import('./settings/load-settings.svelte').then(res => {
      new res.default({
          target: lrs,
          hydrate: true
      });
  })
}

let edp = document.querySelector(`.edit-post`)
if(edp && authenticated) {
  import('../edit-post/edit-post.svelte').then(res => {
      new res.default({
          target: edp,
          hydrate: true
      });
  })
}

let uss = document.querySelector(`.user-sub-spaces`)
if(uss) {
  import('./profile/sub-spaces.svelte').then(res => {
      new res.default({
          target: uss,
          hydrate: true
      });
  })
}

let shp = document.querySelector(`.share-post`)
if(shp && authenticated) {
  import('../share-post/share-post.svelte').then(res => {
      new res.default({
          target: shp,
          hydrate: true
      });
  })
}

let shares = document.querySelectorAll(`.share`)
if(shares && authenticated) {
  import('../post/share/share.svelte').then(res => {
    shares.forEach(share => {
      new res.default({
          target: share,
          props: {
            id: share.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

let menus = document.querySelectorAll(`.post-menu`)
if(menus && authenticated) {
  import('../post/menu/menu.svelte').then(res => {
    menus.forEach(menu => {
      new res.default({
          target: menu,
          props: {
            id: menu.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

let readmores = document.querySelectorAll(`.read-more`)
if(readmores) {
  import('../post/content/read-more/read-more.svelte').then(res => {
    readmores.forEach(readmore => {
      new res.default({
          target: readmore,
          props: {
            id: readmore.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

let rtr = document.querySelectorAll(`.reply-to-reply`)
if(rtr && authenticated) {
  import('../new-post/reply.svelte').then(res => {
    rtr.forEach(reply => {
      new res.default({
          target: reply,
          props: {
            id: reply.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

let nsfw = document.querySelectorAll(`.nsfw`)
if(nsfw) {
  import('../post/nsfw/nsfw.svelte').then(res => {
    nsfw.forEach(item => {
      new res.default({
          target: item,
          props: {
            id: item.dataset.id,
          },
          hydrate: true
      });
    })
  })
}

