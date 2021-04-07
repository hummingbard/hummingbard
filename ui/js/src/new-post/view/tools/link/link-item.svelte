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
      if(res?.is_wikipedia) {
          metadata.is_wikipedia = res.is_wikipedia
      }
      if(res?.sound_cloud_player) {
          metadata.sound_cloud_player = res.sound_cloud_player
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
$: isWikipedia = link?.metadata?.is_wikipedia || false

$: youtube = link?.metadata?.is_youtube || false
$: soundcloud = link?.metadata?.sound_cloud_player?.length > 0
$: vimeo = link?.metadata?.is_vimeo || false
$: imgSrc = youtube ?
`https://img.youtube.com/vi/${link?.metadata?.youtube_id}/mqdefault.jpg` : (vimeo || soundcloud) ? link.metadata?.image : ``


function killMe() {
    dispatch('deleteLink', link.id)
}

onDestroy(() => {
    dispatch('deleteLink', link.id)
})


</script>

<div class="link-item flex">

  {#if youtube || vimeo || soundcloud}
      <div class="vp-i gr-default bg-img"
      style="background-image: url({imgSrc});">
      </div>

  {/if}

  <div class="flex flex-column pa3 flex-one">

      <div class="flex">
          {#if isWikipedia}
            <div class="gr-center ico-20 mr3">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32"><path d="M 30 9.214844 C 30 9.386719 29.863281 9.523438 29.6875 9.523438 L 28.007813 9.523438 L 20.390625 25.738281 C 20.339844 25.847656 20.230469 25.917969 20.113281 25.917969 C 20.109375 25.917969 20.109375 25.917969 20.109375 25.917969 C 19.988281 25.917969 19.882813 25.851563 19.828125 25.746094 L 16.214844 18.578125 L 12.3125 25.757813 C 12.257813 25.859375 12.15625 25.917969 12.03125 25.917969 C 11.914063 25.914063 11.808594 25.847656 11.757813 25.742188 L 4.054688 9.523438 L 2.3125 9.523438 C 2.140625 9.523438 2 9.386719 2 9.214844 L 2 8.390625 C 2 8.21875 2.140625 8.082031 2.3125 8.082031 L 8.523438 8.082031 C 8.695313 8.082031 8.835938 8.21875 8.835938 8.390625 L 8.835938 9.214844 C 8.835938 9.386719 8.695313 9.523438 8.523438 9.523438 L 7.1875 9.523438 L 12.503906 21.785156 L 15.269531 16.617188 L 11.761719 9.527344 L 10.917969 9.527344 C 10.746094 9.527344 10.605469 9.386719 10.605469 9.214844 L 10.605469 8.394531 C 10.605469 8.222656 10.746094 8.082031 10.917969 8.082031 L 15.515625 8.082031 C 15.6875 8.082031 15.824219 8.222656 15.824219 8.394531 L 15.824219 9.214844 C 15.824219 9.386719 15.6875 9.523438 15.515625 9.523438 L 14.703125 9.523438 L 16.722656 13.9375 L 19.125 9.523438 L 17.652344 9.523438 C 17.476563 9.523438 17.339844 9.386719 17.339844 9.214844 L 17.339844 8.394531 C 17.339844 8.222656 17.476563 8.082031 17.652344 8.082031 L 22.117188 8.082031 C 22.289063 8.082031 22.425781 8.222656 22.425781 8.394531 L 22.425781 9.214844 C 22.425781 9.386719 22.289063 9.523438 22.117188 9.523438 L 21.136719 9.523438 L 17.632813 15.894531 L 20.488281 21.769531 L 26 9.523438 L 24.253906 9.523438 C 24.082031 9.523438 23.941406 9.386719 23.941406 9.214844 L 23.941406 8.394531 C 23.941406 8.222656 24.082031 8.082031 24.253906 8.082031 L 29.6875 8.082031 C 29.863281 8.082031 30 8.222656 30 8.394531 Z"/></svg>
            </div>
          {/if}
          <div class="gr-center primary fs-09 flex-one lh-copy">
              {title}
          </div>
          <div class="pointer o-70 hov-op" on:click={killMe}>
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M1 12C1 5.925 5.925 1 12 1s11 4.925 11 11-4.925 11-11 11S1 18.075 1 12zm8.036-4.024a.75.75 0 00-1.06 1.06L10.939 12l-2.963 2.963a.75.75 0 101.06 1.06L12 13.06l2.963 2.964a.75.75 0 001.061-1.06L13.061 12l2.963-2.964a.75.75 0 10-1.06-1.06L12 10.939 9.036 7.976z"></path></svg>
          </div>
      </div>

      {#if descriptionExists}
          <div class="small o-80 mt2 lh-copy" class:clmp-2={!isWikipedia}>
          {#if isWikipedia}
              {@html link.metadata.description}
          {:else}
              {link.metadata.description}
          {/if}
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
