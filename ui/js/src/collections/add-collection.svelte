<script>
import {createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher();

import {addCollection} from './store.js'
import { makeid } from '../utils/utils.js'

let nameInput;
let name = ''
let descriptionInput;
let description = ''

function updateDetails(e) {
}

function kill() {
  dispatch('kill')
}

function save() {
  addCollection({
    id: makeid(16),
    name: name,
    description: description,
  })
  kill()
}

</script>

<div class="gr-center flex flex-column w-100 pa3 brd-btm">

  <div class="">
    <span class=""><strong>New Collection<strong></span>
  </div>

  <div class="mt4">
    <div class="">
      <span class="small"><strong>Name</strong></span>
    </div>
    <div class="mt2 fs-09">
      <input placeholder="My Space Collection" 
       bind:value={name}
        on:keypress={updateDetails}
      autofocus/>
    </div>
  </div>

  <div class="mt3">
    <div class="">
      <span class="small"><strong>Description</strong></span>
    </div>
    <div class="mt2 fs-09">
      <textarea 
        style="height:120px;"
        bind:value={description}
        on:keypress={updateDetails}
        placeholder="Description"
        max-length="500"
        />
    </div>
  </div>


  <div class="mt4 flex">
    <div class="flex-one"></div>

    <div class="">
      <button class="" on:click={save}>Save</button>
      <button class="light" on:click={kill}>Cancel</button>
    </div>
  </div>

</div>
