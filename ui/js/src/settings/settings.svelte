<script>
let confirmDelete = false;

function showConfirmDelete() {
    confirmDelete = true
}

function killConfirmDelete() {
    passwordInput.value = null
    deactivationError = false
    confirmDelete = false
}

$: federated = identity.federated

let form;

let passwordInput;

let deactivationError = false;

let deleting = false;

function confirm() {
    deleting = true
  confirmDeactivation().then((res) => {
    console.log(res)
    if(res?.error) {
        passwordInput.value = null
        deactivationError = true
        passwordInput.focus()
        deleting = false
    }
    if(res?.deactivated) {
        window.location.href = `/logout`
    }
  }).then(() => {
  })
}
async function confirmDeactivation() {
    let endpoint = `/account/delete`
    let data = {
        password: passwordInput.value,
    };
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


function enter(e) {
    if(e.key == 'Enter') {
        confirm()
    }
}

</script>

<div class="flex flex-column lh-copy">

    <div class="brd-btm pa3">
        <div class="">
            <span class="small"><strong>Preferences</strong></span>
        </div>
        <div class="mt3 flex flex-column">
            TODO
        </div>
    </div>


    {#if !federated}
    <div class="pa3">
        <div class="">
            <span class="small"><strong>Delete Account</strong></span>
        </div>
        <div class="mt3 flex flex-column">
            <div class="">
                This action is irreversible. You will not be able to recover your account once it's deleted. 
            </div>
            {#if confirmDelete}
                <div class="mt3 flex flex-column">
                    <span class="small"><strong>Type your password to confirm deactivation</strong></span>

                    <div class="mt3 flex">
                        <div class="flex-one mr3">
                            <input class="" type="password" bind:this={passwordInput}
                            on:keypress={enter}
                          name="password" minlength="8" placeholder="password"
                          autofocus>
                        </div>
                        <div class="gr-center">
                            <button class="warn" on:click={confirm}
                                    disabled={deleting}>{deleting ? 'Deleting...':'Delete'}</button>
                            <button class="light" on:click={killConfirmDelete} disabled={deleting}>Cancel</button>
                        </div>
                    </div>
                    {#if deactivationError}
                        <div class="mt2 fs-09">
                            <strong>Could not delete your account. Perhaps your password is wrong?</strong>
                        </div>
                    {/if}
                </div>
            {:else}
                <div class="mt2">
                    <button class="" on:click={showConfirmDelete}>Delete Account</button>
                </div>
            {/if}
        </div>
    </div>
    {/if}

</div>
