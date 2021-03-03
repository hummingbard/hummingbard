let nav = document.querySelector('.nav-de')
if(nav){
  import('./nav.svelte').then(res => {
    new res.default({
      target: nav,
      hydrate: true,
    })
  })
}
let ies = document.querySelector('.index-explore-spaces')
if(ies){
  import('./index-explore-spaces.svelte').then(res => {
    new res.default({
      target: ies,
      hydrate: true,
    })
  })
}
