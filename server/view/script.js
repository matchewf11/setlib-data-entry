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

// keeping it a buck, i used chatgpt on everything below here
const subjectOptions = {
  data: [
    { id: "avl", value: "avl-tree", label: "AVL Trees" },
    { id: "heap", value: "heap", label: "Heaps" },
    { id: "graph", value: "graph", label: "Graphs" },
  ],
  alg: [
    { id: "sort", value: "sorting", label: "Sorting" },
    { id: "search", value: "searching", label: "Searching" },
    { id: "dp", value: "dynamic-programming", label: "Dynamic Programming" },
  ],
  discrete: [
    { id: "set", value: "set-theory", label: "Set Theory" },
    { id: "logic", value: "logic", label: "Logic" },
    { id: "proof", value: "proof-techniques", label: "Proof Techniques" },
  ],
};

function updateSubjects(section) {
  const container = document.getElementById("subjectOptions");
  container.innerHTML = ""; // Clear existing

  const subjects = subjectOptions[section] || [];

  subjects.forEach((subj) => {
    const radio = document.createElement("input");
    radio.type = "radio";
    radio.id = subj.id;
    radio.name = "subject";
    radio.value = subj.value;

    const label = document.createElement("label");
    label.htmlFor = subj.id;
    label.textContent = subj.label;

    container.appendChild(radio);
    container.appendChild(label);
    container.appendChild(document.createElement("br"));
  });
}

// Attach listeners to section radios
document.addEventListener("DOMContentLoaded", () => {
  const sectionRadios = document.querySelectorAll('input[name="section"]');
  sectionRadios.forEach((radio) => {
    radio.addEventListener("change", () => {
      updateSubjects(radio.value);
    });
  });

  // Optional: initialize to default if one is pre-checked
  const prechecked = Array.from(sectionRadios).find((r) => r.checked);
  if (prechecked) {
    updateSubjects(prechecked.value);
  }
});
