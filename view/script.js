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
  const section =
    document.querySelector('input[name="section"]:checked')?.value || "";
  const difficulty =
    document.querySelector('input[name="difficulty"]:checked')?.value || "";
  const problem = document.getElementById("problem").value;
  const form = document.getElementById("questionForm");
  const previewImage = document.getElementById("previewImage");

  if (!section || !difficulty || !problem) return;

  fetch("/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ section, difficulty, problem }),
  })
    .then(async (response) => {
      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.message || "Submit failed");
      }
      return response.json();
    })
    .then(() => {
      form.reset();
      previewImage.style.display = "none";
    })
    .catch((error) => {
      console.error("Submit error:", error.message);
    });
}
