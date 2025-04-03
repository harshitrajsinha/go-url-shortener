document.addEventListener("DOMContentLoaded", function () {
  document
    .querySelector(".shorten-url-btn")
    .addEventListener("click", function () {
      const originalURL = document.getElementById("original-url").value;
      document.getElementById("original-url").value = "";
      const requestURL = window.location.href + "api/v1/shorten";
      const shortenedURLBlock = document.querySelector(".shortened-url");

      if (originalURL) {
        shortenedURLBlock.style.color = "white";
        shortenedURLBlock.textContent = "Loading...";
        const data = {
          url: originalURL,
        };

        fetch(requestURL, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        })
          .then((response) => {
            if (!response.ok) {
              shortenedURLBlock.textContent = "";
              shortenedURLBlock.style.Color = "red";
              if (response.status === 409) {
                shortenedURLBlock.textContent = "URL already shortened";
              } else {
                shortenedURLBlock.textContent = "Error shortening this URL";
              }

              return response.text().then((errorText) => {
                return Promise.reject([
                  response.status,
                  `Error: ${response.statusText} - ${errorText}`,
                ]);
              });
            }
            return response.json();
          })
          .then((data) => {
            if (Number(data["code"]) === 201) {
              shortenedURLBlock.style.color = "#7deb7d";
              const shortenedURL = data["data"][0]["shortened-url"];
              shortenedURLBlock.textContent = shortenedURL;
            }
          })
          .catch((error) => {
            if (error[0] === 409) {
              shortenedURLBlock.style.color = "red";
              shortenedURLBlock.textContent = error[1]
                .split("[")[1]
                .split("]")[0];
            }
            alert(error[1]);
          });
      }
    });
});
