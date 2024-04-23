<script>
    import '../app.css'
    import { Toaster } from 'svelte-sonner';
</script>

<Toaster richColors position="bottom-left" duration={7000} />
<slot />