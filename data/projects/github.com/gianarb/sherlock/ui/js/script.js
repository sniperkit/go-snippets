const fetchTrace = (e) => {
  e.preventDefault();
  const id = $('.search-input').val();

  return fetch(`/api/trace/${id}`)
    .then(res => {
      if (res.status !== 200) {
        console.log(res);
      }
      return res;
    })
    .then(res => res.json())
    .then(res => console.log({ res }))
    .catch(error => console.log(error));
}

$('.submit-btn').on('click', fetchTrace);
