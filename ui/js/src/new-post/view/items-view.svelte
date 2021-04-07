<script>
import LinkItems from './tools/link/items.svelte'
import ImageItems from './tools/image/items.svelte'
import AttachmentItems from './tools/attachment/items.svelte'
import ReviewView from './tools/review/view.svelte'
import {post} from '../store.js'

export let store;

</script>

{#if $post.active || store.active}

    {#if store?.links?.length > 0}
        <LinkItems 
        store={store} 
        on:deleteLink
        on:updateLinkMetadata/>
    {/if}

    {#if store?.images?.length > 0}
        <ImageItems
        store={store} 
        on:deleteImage
        on:updateImageURL
        on:updateImageMetadata/>
    {/if}

    {#if store?.attachments && store?.attachments?.length > 0}
        <AttachmentItems
        store={store} 
        on:deleteAttachment
        on:updateAttachmentURL/>
    {/if}

    {#if store?.review?.enabled }
        <ReviewView
        store={store} 
        on:updateReview/>
    {/if}

{/if}

