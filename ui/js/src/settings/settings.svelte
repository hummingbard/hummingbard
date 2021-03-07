<script>
import Preferences from './preferences.svelte'
import Delete from './delete.svelte'

$: federated = identity.federated
</script>

<div class="flex flex-column lh-copy">

    <Preferences />


    {#if !federated}
        <Delete/>
    {/if}

</div>
