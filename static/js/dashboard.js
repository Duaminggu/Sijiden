function userCounter() {
    return {
        count: '...',
        async load() {
            try {
                const response = await fetch('/ajax/sijiden/users/count');
                const result = await response.json();
                this.count = result.count || 0;
            } catch (error) {
                console.error('Failed to load user count:', error);
                this.count = 'ERR';
            }
        }
    }
}