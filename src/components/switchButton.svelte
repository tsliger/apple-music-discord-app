<script lang="ts">
    import { enable, isEnabled, disable } from "@tauri-apps/plugin-autostart";
    import { onMount } from "svelte";

    let enabledState = true;
    onMount(async () => {
        enabledState = await isEnabled();
    });

    const setAutostart = async (event: MouseEvent) => {
        event.preventDefault();

        if ((await isEnabled()) == false) {
            await enable();
        } else {
            await disable();
        }

        // update frontend state
        enabledState = await isEnabled();
    };
</script>

<label class="switch">
    <input type="checkbox" onclick={setAutostart} bind:checked={enabledState} />
    <span class="slider round"></span>
</label>

<style>
    .switch {
        position: relative;
        display: inline-block;
        width: 36px;
        height: 20px;
    }

    .switch input {
        opacity: 0;
        width: 0;
        height: 0;
    }

    .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: #ccc;
        -webkit-transition: 0.4s;
        transition: 0.4s;
    }

    .slider:before {
        position: absolute;
        content: "";
        height: 12px;
        aspect-ratio: 4/4;
        left: 4px;
        bottom: 4px;
        background-color: white;
        -webkit-transition: 0.4s;
        transition: 0.4s;
    }

    input:checked + .slider {
        background-color: #2196f3;
    }

    input:focus + .slider {
        box-shadow: 0 0 1px #2196f3;
    }

    input:checked + .slider:before {
        -webkit-transform: translateX(16px);
        -ms-transform: translateX(16px);
        transform: translateX(16px);
    }

    .slider.round {
        border-radius: 34px;
    }

    .slider.round:before {
        border-radius: 50%;
    }
</style>
