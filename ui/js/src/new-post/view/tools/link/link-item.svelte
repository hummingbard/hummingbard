<script>
import { createEventDispatcher } from 'svelte';
const dispatch = createEventDispatcher();

export let store;

import {onMount, onDestroy} from 'svelte'

export let link;

let adding = true;

onMount(() => {
  fetchMetadata().then((res) => {
    console.log(res)
    let metadata = {}
      if(res?.title) {
          metadata.title = res.title
      }
      if(res?.author) {
          metadata.author = res.author
      }
      if(res?.description) {
          metadata.description = res.description
      }
      if(res?.image) {
          metadata.image = res.image
      }
      if(res?.is_youtube) {
          metadata.is_youtube = res.is_youtube
      }
      if(res?.youtube_id) {
          metadata.youtube_id = res.youtube_id
      }
      if(res?.is_vimeo) {
          metadata.is_vimeo = res.is_vimeo
      }
      if(res?.vimeo_id) {
          metadata.vimeo_id = res.vimeo_id
      }
      dispatch('updateLinkMetadata', {
        id: link.id,
        metadata: metadata,
      })
      adding = false
  }).then(() => {
  })

})

async function fetchMetadata() {
    let endpoint = `/link/metadata`

    let data = {
        href: link.href,
    }

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: JSON.stringify(data),
    headers:{
      'Authorization': identity.access_token,
        'Content-Type': 'application/json'
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

$: title = link?.metadata?.title?.length > 0 ? link.metadata.title : link.href
$: titleExists = link?.metadata?.title?.length > 0
$: descriptionExists = link?.metadata?.description?.length > 0

$: youtube = link?.metadata?.is_youtube || false
$: vimeo = link?.metadata?.is_vimeo || false
$: imgSrc = youtube ? `https://img.youtube.com/vi/${link?.metadata?.youtube_id}/mqdefault.jpg` : vimeo ? link.metadata?.image : ``


function killMe() {
    dispatch('deleteLink', link.id)
}

onDestroy(() => {
    dispatch('deleteLink', link.id)
})


</script>

<div class="link-item flex">

  {#if youtube || vimeo}
      <div class="vp-i gr-default bg-img"
      style="background-image: url({imgSrc});">
      </div>

  {/if}

  <div class="flex flex-column pa3 flex-one">

      <div class="flex">
          <div class="primary fs-09 flex-one lh-copy">
              {title}
          </div>
          <div class="pointer o-70 hov-op" on:click={killMe}>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M1 12C1 5.925 5.925 1 12 1s11 4.925 11 11-4.925 11-11 11S1 18.075 1 12zm8.036-4.024a.75.75 0 00-1.06 1.06L10.939 12l-2.963 2.963a.75.75 0 101.06 1.06L12 13.06l2.963 2.964a.75.75 0 001.061-1.06L13.061 12l2.963-2.964a.75.75 0 10-1.06-1.06L12 10.939 9.036 7.976z"></path></svg>
          </div>
      </div>

      {#if descriptionExists}
      <div class="small o-80 mt2 clmp-2 lh-copy">
          {link.metadata.description}
      </div>
      {/if}

      {#if titleExists}
      <div class="small o-80 mt2 cmpl-2">
          {link.href}
      </div>
      {/if}


      {#if adding}
          <div class="fs-09 mt3">
              <em>Fetching Link Metadata...</em>
          </div>
      {/if}
  </div>

</div>

<style>
.link-item {
    border: 1px solid var(--primary-light-gray);
    border-radius: 7px;
    margin-bottom:1rem;
}
</style>
