<script>
import Avatar from '../../../avatar/avatar.svelte'
import { settings } from '../store.js'

function updateHeader(e) {
  console.log("lol", e.detail)
  $settings.appearance.header = e.detail
}

function removeHeader(e) {
  $settings.appearance.header = null
}


</script>

<div class="flex flex-column flex-one appearance">


  <div class="flex flex-column">

    <div class="flex flex-column">
      <div class="">
        <span class="small"><strong>Header Image</strong></span>
      </div>
      <div class="mt2 h-100">
        <Avatar 
          header={true}
          on:uploaded={updateHeader}
          on:removed={removeHeader}
          avatar={$settings.appearance.header}/>
      </div>
    </div>

    <div class="mt3 flex flex-column h-100">
      <div class="">
        <span class="small"><strong>Custom CSS</strong></span>
      </div>
      <div class="mt2 h-100">
        <textarea
          style="height:200px;"
          placeholder="Add custom CSS here."
          bind:value={$settings.appearance.css}
        autofocus></textarea>
      </div>
    </div>

  </div>

</div>

<style>

input, textarea {
  width: 100%;
}
</style>
