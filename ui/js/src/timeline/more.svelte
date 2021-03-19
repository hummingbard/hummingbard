<script>
import { onMount } from 'svelte'
import { addPosts } from './store.js'

export let type;

let last = window?.timeline?.end;

let posts = [];

let loading = false;
let noMore = false;


let Post;
let postLoaded = false;
function loadPost() {
    import('../post/post.svelte').then(res => {
        Post = res.default
        postLoaded = true
    })
}

onMount(() => {
    clientHeight = document.body.clientHeight
    observer = new IntersectionObserver(callback, options);
    observer.observe(obs);
})

let feed = window.timeline?.feed

async function fetchFollowing() {
    let endpoint = `/messages/fetch`

    let data = {
        id: window.timeline.room_id,
        end: last,
    }

    if(window.timeline?.showAll) {
        data.all = true
    }

    if(window.timeline.userFeed) {
        console.log(feed)
        endpoint = `/feed/fetch`
        data = {
            feed: feed,
        }
    }

    if(window.timeline?.public) {
        data.public = true
    }


  if(window?.timeline?.permalink) {
    data.id = window.timeline.thread_in_room_id
      data.permalink = true
  }

    let options = {
        method: 'POST', // or 'PUT'
        body: JSON.stringify(data),
        headers:{
            'Content-Type': 'application/json'
        }
    }



    if(authenticated && identity?.access_token) {
        options.headers['Authorization'] = identity.access_token
    }
    console.log(options)

    let resp = await fetch(endpoint, options)
    const ret = await resp.json()
    return Promise.resolve(ret)
}

function loadMore() {
     loading = true
  fetchFollowing().then((res) => {
    console.log("Fetched posts: ",res)
    if(!postLoaded) {
        loadPost()
    }
      if(res?.last_event) {
          last = res.last_event
      }
      if(res?.posts?.length > 0) {
          addPosts(res.posts)
          res.posts.forEach(post => {
              if(post.type == `com.hummingbard.post` || post.type == `m.room.message` && !post.redacted) {
                posts = [...posts, post];
              }
          })
      }
      if(res?.last_event == `t0_0` || res?.posts?.length == 0) {
          noMore = true
      }
      if(res?.feed_items) {
          feed = res.feed_items
      }
  }).then(() => {
     loading = false
  })
}


let observer;
let obs;
let clientHeight = document.body.clientHeight;
let options = {
    root: obs,
    rootMargin: `${clientHeight / 2.7}px`,
    threshold: 0.5,
}


function callback(entries, observer) {
  entries.forEach(entry => {
    if(entry.isIntersecting) {
      loading = true
      loadMore()
    }
  })
}


</script>


<div class="more-posts flex flex-column">
  <div class="more-post-items">
    {#if postLoaded}

        {#each posts as post (post.event_id)}
            <svelte:component this={Post} post={post}/>
        {/each}

    {/if}
  </div>
    <div class="load-more tc pv5">
        {#if !noMore}
            {#if loading}
                <div class="lds-ring"><div></div><div></div><div></div><div></div></div>
            {:else}
                <button class="" on:click={loadMore}>Load More</button>
            {/if}
        {:else}
            <span class="small bold">No More Posts</span>
        {/if}
    </div>
  <div bind:this={obs}>
  </div>
</div>

