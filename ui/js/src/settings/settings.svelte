<script>
import Preferences from './preferences.svelte'
import Password from './password.svelte'
import Delete from './delete.svelte'
import Email from './email.svelte'

$: federated = identity.federated

</script>

<div class="flex flex-column lh-copy">

    <Preferences />

    <Email />

    <Password/>


    {#if !federated}
        <Delete/>
    {/if}

</div>
