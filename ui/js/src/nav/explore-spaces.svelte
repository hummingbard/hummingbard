<script>
import {onMount, createEventDispatcher} from 'svelte'
import {fade} from 'svelte/transition'
import SpaceItem from './explore-space-item.svelte'
import {debounce} from '../utils/utils.js'
const dispatch = createEventDispatcher();
let fetched = false;fetch
let spaces = [];
let results = [];



function filtered() {
    return spaces.filter(x => x.room_alias.toLowerCase().includes(query) || x.topic.toLowerCase().includes(query) || x.name.toLowerCase().includes(query))
}

onMount(() => {
    fetchSpaces().then(res => {
        console.log(res)
        if(res?.spaces?.length > 0) {
            spaces.push(...res?.spaces)
            fetched = true
        }
    }).then(() => {
        if(fetched && spaces.length > 0) {
            observer = new IntersectionObserver(callback, options);
            observer.observe(obs);
        }
    })
})

let observer;
let obs;
let options = {
    rootMargin: `0px`,
    threshold: 0.00,
}


let noMore = false;

function loadMore() {
    if(noMore) {
        return
    }
    if(spaces.length < 13) {
        noMore = true
        return
    }
    fetchSpaces().then(res => {
        console.log(res)
        if(res?.spaces?.length > 0) {
            res?.spaces.forEach(space => {
                if(spaces.findIndex(x => x.room_id == space.room_id) == -1) {
                    spaces = [...spaces, space]
                }
            })
            fetched = true
        } else if(res?.spaces.length == 0) {
            noMore = true
        }
    }).then(() => {
    })
}

function callback(entries, observer) {
  entries.forEach(entry => {
    if(entry.isIntersecting) {
      loadMore()
    }
  })
}

function kill() {
    dispatch('kill', true)
}

async function fetchSpaces() {
    let endpoint = `/spaces/public`

    let data = {
        offset: 0,
    }

    if(spaces.length > 0) {
        data.offset = spaces.length
    }

    let resp = await fetch(endpoint, {
    method: 'POST', 
    body: JSON.stringify(data),
    headers:{
        'Content-Type': 'application/json'
    }
    })
    if (resp.ok) { 
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

async function querySpace() {
    let endpoint = `/spaces/query`

    let data = {
        query: query,
    }

    let resp = await fetch(endpoint, {
    method: 'POST', 
    body: JSON.stringify(data),
    headers:{
        'Authorization': identity.access_token,
        'Content-Type': 'application/json'
    }
    })
    if (resp.ok) { 
    } else {
      alert("HTTP-Error: " + resp.status);
    }
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let searching = false
let searchInput;
let query = '';

let fetchingResults = false;

$: if(query.length > 0) {
    searching = true
    fetchingResults = true
} else if(query.length == 0) {
    searching = false
    results = []
    results = results
}

function cancelSearch() {
    searching = false
    query = ''
    searchInput.value = ''
    results = []
    results = results
}

function updated(e) {
    debounce(() => {
        querySpace().then(res => {
            console.log(res)
            if(!searching) {
                return
            }
            if(res?.spaces?.length > 0) {
                results = res?.spaces
                results =results
            }
        }).then(() => {
            fetchingResults = false
        })
    }, 500)
}

</script>


<div class="modal-container main pa3" transition:fade="{{duration: 73}}">

  <div class="modal-inner-np start flex flex-column " >

    <div class="pa3 flex flex-column">

        <div class="">
            <span class=""><strong>Explore Spaces</strong></span>
        </div>

        <div class="mt3 lh-copy">
            <span class="">If you can't find the space you're looking for, <a
                  href="/create"><span class="primary">Create a new
                      space.</span></a></span>
        </div>

        <div class="mt3 fs-09 relative">
            <input bind:this={searchInput} bind:value={query}
                   on:keydown={updated}
            placeholder="Filter Spaces"/>
            {#if query.length > 0}
                <div class="absolute disc pointer gr-default mr2" on:click={cancelSearch}>
                    <svg class="gr-center" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
                </div>
            {/if}
        </div>
        {#if !searching}
            {#if !fetched}
                <div class="gr-default w-100 pt4">
                    <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
                </div>
            {:else}
                {#if spaces.length > 0}
                    <div class="flex flex-column mt3 mx" >
                        <div class="ovfl-y scrl">
                        {#each spaces as space (space.room_id)}
                            <SpaceItem 
                            space={space} />
                        {/each}
                              <div bind:this={obs}>
                              </div>
                              <div class="mt3">
                                  {#if !noMore && query.length == 0}
                                      <button class="" on:click={loadMore}>Load More</button>
                                  {/if}
                              </div>
                        </div>
                    </div>
                {:else}
                    <div class="pa4 tc">
                        <span>No spaces</span>
                    </div>
                {/if}
            {/if}

        {:else}
                {#if results?.length > 0}
                    <div class="flex flex-column mt3 mx" >
                        <div class="ovfl-y scrl">
                        {#each results as space (space.room_id)}
                            <SpaceItem 
                            space={space} />
                        {/each}
                        </div>
                    </div>
                {:else}
                    {#if fetchingResults}
                        <div class="gr-default w-100 pt4">
                            <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
                        </div>
                    {:else}
                        <div class="pa4 tc">
                            <span>Didn't find anything...</span>
                        </div>
                    {/if}
                {/if}
        {/if}




    </div>


  </div>

  <div class="mask absolute" on:click={kill}></div>

</div>

<style>
.modal-container {
    top: 0;
    left: 0;
    position: fixed;
    width: 100%;
    height: 100%;
    z-index: 49999;
    display: grid;
    grid-template-columns: auto;
    grid-template-rows: auto;
}

.modal-inner-np {
    justify-self: center;
    display: grid;
    grid-template-columns: auto;
    grid-template-rows: auto;
    z-index: 5000;

    background-color: var(--m-bg);
    -webkit-box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    -moz-box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    box-shadow: 0px 4px 24px -1px rgba(0,0,0,0.05);
    border-radius: 17px;
    transition: 0.1s;
    word-break: break-word;
}

.start {
    align-self: center;
    width: 680px;
    height: 70vh;
}

  .mx {
      height: 52vh;
  }

@media screen and (max-width: 680px) {
  .start {
    width: 100%;
    height: 100vh;
  }
}

.mask {
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background: var(--mask)
}

.disc {
    right:0;
    top: 0;
    height: 100%;
}

input, textarea {
  width: 100%;
}

</style>
