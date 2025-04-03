document.addEventListener('alpine:init', () => {
    Alpine.directive('page-active', (el, { expression }, { evaluate }) => {
        const page = evaluate(expression);

        const highlightLinks = () => {
            const links = el.querySelectorAll('[data-page]');
            links.forEach(link => {
                link.classList.toggle('active', link.dataset.page === page);
            });
        };

        // langsung coba dulu, siapa tahu sudah ready
        highlightLinks();

        // siapkan observer kalau isi belum muncul
        const observer = new MutationObserver(() => {
            highlightLinks();
        });

        observer.observe(el, { childList: true, subtree: true });
    });
});