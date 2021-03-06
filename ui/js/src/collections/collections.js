let col = document.querySelector('.collections')
if(col) {
  import('./collections.svelte').then(res => {
    new res.default({
      target: col,
      hydrate: true,
    })
  })
}
