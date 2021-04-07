<script>
import {onMount} from 'svelte'
import Star from './star.svelte'
import { createEventDispatcher } from 'svelte';
const dispatch = createEventDispatcher();

export let store;

let titleInput;
let contentInput;

let rating = store.review?.rating

function updateReview() {
    dispatch('updateReview', {
        title: titleInput.value,
        content: contentInput.value,
        rating: rating,
    })
}

function rate(e) {
    rating = e.detail
    updateReview()
}

</script>

<div class="flex flex-column">

    <div class="">
      <input 
        bind:this={titleInput}
        on:input={updateReview}
      placeholder="Lord of the Rings" autofocus/>
    </div>

    <div class="mt3">
      <textarea 
        bind:this={contentInput}
        on:input={updateReview}
         maxlength="1000"
      placeholder="Promising author. Needs to work on his prose."></textarea>
    </div>

    <div class="mt3 mb3">
        <div class="flex">
            {#each Array(7) as _, i}
                <Star index={i} store={store} on:rate={rate}/>
            {/each}

        </div>
    </div>

</div>

<style>
input, textarea {
    width: 100%;
}

textarea {
    height: 120px;
}
</style>
