document.addEventListener('alpine:init', () => {
    const cache = new Map();

    Alpine.directive('component', (el, { modifiers, expression }) => {
        const loadComponent = () => {
            if (modifiers.includes('once') && cache.has(expression)) {
                el.innerHTML = cache.get(expression);
                Alpine.initTree(el);
                return;
            }

            fetch(expression)
                .then(res => res.text())
                .then(html => {
                    el.innerHTML = html;
                    Alpine.initTree(el);
                    if (modifiers.includes('once')) {
                        cache.set(expression, html);
                    }
                });
        };

        if (modifiers.includes('lazy')) {
            const observer = new IntersectionObserver(([entry], observer) => {
                if (entry.isIntersecting) {
                    loadComponent();
                    observer.unobserve(el);
                }
            });
            observer.observe(el);
        } else {
            loadComponent();
        }
    });
});
