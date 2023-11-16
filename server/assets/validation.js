let loginModal = document.getElementById('loginModal');
let statusElm = document.getElementById('status');

function openModal() {
  loginModal.style.display = 'block';
}

function closeModal() {
  loginModal.style.display = 'none';
}

function saveUUID(uuid) {
  localStorage.setItem('verificationUUID', uuid);
}

function validateLogin() {
  let username = document.getElementById('username').value;
  let password = document.getElementById('password').value;

  data = {
    username: username,
    password: password,
  };
  const jsonString = JSON.stringify(data);

  fetch('/auth/validate', {
    method: 'POST',
    body: jsonString,
  })
    .then((response) => response.json())
    .then((data) => {
      // Handle the response from the server
      console.log(data);

      if (data.status === 'success') {
        // Save the received UUID
        saveUUID(data.uuid);
        window.location.href = '/view-data';
      } else if (data.status === 'Invalid credentials') {
        console.log('status');

        statusElm.textContent = 'Invalid credentials';
      }
    })
    .catch((error) => {
      // Handle errors
      console.error('Error:', error);
    });
}

// Add event listener to the form for the submit event
document
  .getElementById('loginForm')
  .addEventListener('submit', function (event) {
    event.preventDefault();
    validateLogin(); // Call your validation function here
  });

openModal();
