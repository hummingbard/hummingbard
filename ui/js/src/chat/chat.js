let chat = document.querySelector('.chat')
console.log("hehe")
if(chat) {
  import('./chat.svelte').then(res => {
    new res.default({
      target: chat,
      hydrate: true,
    })
  })
}
