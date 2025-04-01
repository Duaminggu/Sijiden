function loginForm() {
    return {
        email: '',
        password: '',
        error: '',
        async submit() {
            const csrfToken = document.querySelector('meta[name="csrf-token"]').getAttribute('content');
            console.log(this.email, this.password)
            if (!this.email || !this.password) {
                this.error = 'Email and password are required.';
                return;
            }

            const data = {
                email: this.email,
                password: this.password
            };

            const res = await fetch('/ajax/sijiden/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRF-Token': csrfToken
                },
                body: JSON.stringify(data),
                credentials: 'include'
            });

            const json = await res.json();

            console.log(json)

            if (res.ok) {
                window.location.href = json.redirectTo; // redirect ke halaman todos
            } else {
                alert(json.error);
            }

        }
    }
}

// const forms = xa2csp859r.querySelector(".forms"),
//     pwShowHide = xa2csp859r.querySelectorAll(".eye-icon"),
//     links = xa2csp859r.querySelectorAll(".link");
// // Add click event listener to each eye icon for toggling password visibility
// pwShowHide.forEach(eyeIcon => {
//     eyeIcon.addEventListener("click", () => {
//         let pwFields = eyeIcon.parentElement.parentElement.querySelectorAll(".password");
//         pwFields.forEach(password => {
//             if (password.type === "password") { // If password is hidden
//                 password.type = "text"; // Show password
//                 eyeIcon.classList.replace("bx-hide", "bx-show"); // Change icon to show state
//                 return;
//             }
//             password.type = "password"; // Hide password
//             eyeIcon.classList.replace("bx-show", "bx-hide"); // Change icon to hide state
//         });
//     });
// });
// // Add click event listener to each link to toggle between forms
// links.forEach(link => {
//     link.addEventListener("click", e => {
//         e.preventDefault(); // Prevent default link behavior
//         forms.classList.toggle("show-signup");
//     });
// });

