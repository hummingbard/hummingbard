<script>
import {createEventDispatcher} from 'svelte';
const dispatch = createEventDispatcher();

export let store;
import NewFeaturedImage from '../../../image/new-featured-image.svelte'

function updateFeaturedImage(e) {
  console.log(e.detail)
  dispatch('updateFeaturedImage', {
    content_uri: e.detail.content_uri,
    height: e.detail.height,
    width: e.detail.width,
    caption: '',
  })
}

function remove(e) {
  dispatch('removeFeaturedImage')
}

let subtitle;
let description;
let canonical;
function updateDetails(e) {
  dispatch('updateDetails', {
    subtitle: subtitle.value,
    description: description.value,
    canonical_link: canonical.value,
  })
}

</script>

<div class="flex flex-column">

  <div class="flex flex-column pa3">

    <div class="flex flex-column">
      <div class="flex">
        <span class="small"><strong>Subtitle</strong></span>
      </div>
      <div class="flex mt2 small">
        <input 
          bind:this={subtitle}
          on:keypress={updateDetails}
          value={store?.article?.subtitle}
          placeholder="Subtitle"
          max-length="200"
          autofocus/>
      </div>
    </div>

    <div class="flex flex-column mt3">
      <div class="flex">
        <span class="small"><strong>Description</strong></span>
      </div>
      <div class="flex mt2 small">
        <textarea 
          style="height:120px;"
          bind:this={description}
          on:keypress={updateDetails}
          value={store?.article?.description}
          placeholder="Description"
          max-length="500"
          />
      </div>
    </div>

    <div class="flex flex-column mt4">
      <div class="flex">
        <span class="small"><strong>Featured Image</strong></span>
      </div>
      <div class="flex mt2 small">
        <NewFeaturedImage 
          image={store?.article?.featured_image?.content_uri}
          on:uploaded={updateFeaturedImage}
          on:removed={remove}/>
      </div>
    </div>

    <div class="flex flex-column mt4">
      <div class="flex">
        <span class="small"><strong>Canonical Link</strong></span>
      </div>
      <div class="flex mt2 small">
        <input 
          bind:this={canonical}
          on:keypress={updateDetails}
          value={store?.article?.canonical_link}
          placeholder="Canonical Link"/>
      </div>
    </div>

  </div>

  <div class="flex">
  </div>

</div>
<style>

input, textarea {
  width: 100%;
}
</style>
