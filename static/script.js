document.addEventListener("DOMContentLoaded", function () {
    const tasks = document.querySelectorAll(".task-item");

    tasks.forEach(task => {
        const dueDateString = task.getAttribute("data-due-date");
        if (!dueDateString) return;

        const dueDate = new Date(dueDateString);
        const now = new Date();
        const timeRemaining = (dueDate - now) / (1000 * 60 * 60 * 24); // Convert to days

        if (timeRemaining <= 3) {
            // Make task blink
            setInterval(() => {
                task.classList.toggle("blink");
            }, 5000);
        }
    });
});
