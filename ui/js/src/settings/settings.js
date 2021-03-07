let el = document.querySelector('.settings')
if(el) {
  import('./settings.svelte').then(res => {
    new res.default({
      target: el,
      hydrate: true,
    })
  })
}
