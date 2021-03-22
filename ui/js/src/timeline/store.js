import { writable, derived } from 'svelte/store';

  let posts = [];

  const storePosts = writable(posts);

  let initializePosts = () => {
    if(window.timeline?.initialPosts) {
      window.timeline.initialPosts.forEach(post => {
        post.hydrated = true
        posts.push(post)
        if(post.shared_post) {
          posts.push(post.shared_post)
        }
      })
    }
  }


  let addPosts = (p) => {
    p.forEach(post => {
        posts.push(post)
    })
  }

  let addPost = (post) => {
    storePosts.update(p => {
      post.adding_postscript = false
      p.push(post)
      return p
    })
  }


  let postExists = (id) => {
    let ind = posts.findIndex(x => x.id === id)
    if(ind != -1) {
      return true
    }
    return false
  }


  let addPostScript = (thing) => {
    storePosts.update(p => {
      let ind = p.findIndex(x => x.thing_id === thing)
      if(ind != -1) {
        p[ind].adding_postscript = true
      }
      return p
    })
  }


  let getPostById = (id) => {
    return posts.filter(post => post.event_id === id)[0]
  }

  let getPostByShortId = (id) => {
    return posts.filter(post => post.short_id === id)[0]
  }

  let getPostByShortlink = (shortlink) => {
    return posts.filter(post => post.shortlink === shortlink)[0]
  }

  let getPostByThingId = (thing_id) => {
    return posts.filter(post => post.thing_id === thing_id)[0]
  }


  let getPostReactions = (id) => {
    return posts.filter(post => post.content?.m_relates_to?.event_id === id)
  }

  let addingPS = (thing_id) => {
    return posts.filter(post => post.thing_id === thing_id)[0].adding_postscript
  }


  let queryState = (thing_id) => {
    let ind = posts.findIndex(x => x.thing_id === thing_id)
    if(ind != -1) {
      return posts[ind].initial
    }
    return false
  }


  let updateContent = (thing, item) => {
    storePosts.update(p => {
      let ind = p.findIndex(x => x.thing_id === thing)
      if(ind != -1) {
        p[ind].content = item.content
        p[ind].content_data = item.contentData
      }
      return p
    })
  }


  export {
    storePosts,
    initializePosts,
    addPosts,
    addPost,
    getPostById,
    getPostByShortlink,
    getPostByShortId,
    getPostByThingId,
    getPostReactions,
    queryState,
    postExists,
  }


function createAddedPosts () {
  let added = []
  const { subscribe, set, update } = writable(added);

  let add = (post) => {
    update(p =>  {
      p.unshift(post)
      return p
    })
  }

  return {
    subscribe,
    set,
    add,
  };

}
export const addedPosts = createAddedPosts();



