let chat = document.querySelector('.chat')
if(chat) {
  import('./chat.svelte').then(res => {
    new res.default({
      target: chat,
      hydrate: true,
    })
  })
}
