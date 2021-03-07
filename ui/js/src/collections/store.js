import { writable, derived } from 'svelte/store';

  let collections = [];

  const storeCollections = writable(collections);

  let initializeCollections = () => {
    console.log("collections is", collections)
    storeCollections.update(p => {
      if(window.timeline?.collections) {
        console.log("adding collections")
        window.timeline.collections.forEach(collection => {
          p.push(collection)
        })
      }
      return p
    })
  }


  let addCollection = (collection) => {
    storeCollections.update(p => {
      p.push(collection)
      return p
    })
  }


  let pathExists = (id) => {
    let ind = collections.findIndex(x => x.path === path)
    if(ind != -1) {
      return true
    }
    return false
  }

  export {
    storeCollections,
    initializeCollections,
    addCollection,
    pathExists,
  }


