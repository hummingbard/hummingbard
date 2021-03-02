<script>
import {onMount} from 'svelte'
import {fade} from 'svelte/transition'
import SpaceItem from './explore-space-item.svelte'
import { createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher();
let fetched = false;fetch
let spaces = [];


function filtered() {
    return spaces.filter(x => x.room_alias.toLowerCase().includes(query) || x.topic.toLowerCase().includes(query) || x.name.toLowerCase().includes(query))
}

$: items = query.length > 0 ? filtered() : spaces

onMount(() => {
    fetchSpaces().then(res => {
        console.log(res)
        if(res?.spaces?.length > 0) {
            spaces.push(...res?.spaces)
            fetched = true
        }
    }).then(() => {
        observer = new IntersectionObserver(callback, options);
        observer.observe(obs);
    })
})

let observer;
let obs;
let options = {
    rootMargin: `0px`,
    threshold: 0.00,
}


function loadMore() {
    fetchSpaces().then(res => {
        console.log(res)
        if(res?.spaces?.length > 0) {
            res?.spaces.forEach(space => {
                if(spaces.findIndex(x => x.room_id == space.room_id) == -1) {
                    spaces = [...spaces, space]
                }
            })
            fetched = true
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

let searching = false
let searchInput;
let query = '';


</script>


<div class="modal-container main ph3" transition:fade="{{duration: 73}}">

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

        <div class="mt3 fs-09">
            <input bind:this={searchInput} bind:value={query}
            placeholder="Filter Spaces"/>
        </div>

            {#if !fetched}
                <div class="gr-default w-100 pt4">
                    <div class="lds-ring gr-center"><div></div><div></div><div></div><div></div></div>
                </div>
            {:else}
                {#if spaces.length > 0}
                    <div class="flex flex-column mt3 mx" >
                        <div class="ovfl-y scrl">
                        {#each items as space (space.room_id)}
                            <SpaceItem space={space} />
                        {/each}
                      <div bind:this={obs}>
                      </div>
                        </div>
                    </div>
                {:else}
                    <div class="pa4 tc">
                        <span>No spaces</span>
                    </div>
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
      max-height: 53vh;
  }

@media screen and (max-width: 680px) {
  .start {
    width: 100%;
  }
}

.mask {
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    background: var(--mask)
}
</style>
