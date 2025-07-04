function previewPost() {
  const problem = document.getElementById("problem").value;
  const previewImage = document.getElementById("previewImage");

  fetch("/preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ problem }),
  })
    .then((response) => {
      if (!response.ok) throw new Error("Failed to load preview image");
      return response.blob();
    })
    .then((blob) => {
      previewImage.src = URL.createObjectURL(blob);
      previewImage.style.display = "block";
    })
    .catch((error) => {
      console.error("Preview error:", error.message);
    });
}

function submitPost() {
  const form = document.getElementById("questionForm");
  const previewImage = document.getElementById("previewImage");
  const usernameInput = document.getElementById("username");

  const data = {
    username: usernameInput.value.trim(),
    type: document.querySelector('input[name="type"]:checked')?.value || "",
    section: document.querySelector('input[name="section"]:checked')?.value || "",
    subject: document.querySelector('input[name="subject"]:checked')?.value || "",
    difficulty: document.querySelector('input[name="difficulty"]:checked')?.value || "",
    problem: document.getElementById("problem").value.trim(),
  };

  if (Object.values(data).some(v => !v)) {
    alert("Please fill out all fields.");
    return;
  }

  fetch("/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  })
    .then(res => res.ok ? res.json() : res.text().then(msg => Promise.reject(new Error(msg || "Submit failed"))))
    .then(() => {
      const name = data.username;
      form.reset();
      usernameInput.value = name;
      previewImage.src = "";
      previewImage.style.display = "none";
    })
    .catch(err => console.error("Submit error:", err.message));
}
