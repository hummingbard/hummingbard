import { writable, derived } from 'svelte/store';

  let collections = [];

  const storeCollections = writable(collections);

  let initializeCollections = () => {
    if(window.timeline?.collections) {
      window.timeline.collections.forEach(collection => {
        collection.hydrated = true
        collections.push(collection)
      })
    }
  }


  let addCollection = (collection) => {
    console.log("adding", collection)
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


