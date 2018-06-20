
export function loadSnippets () {
  return (dispatch) => {
    fetch('http://localhost:5100/state', {
      credentials: 'include'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot load state:', e)
    })
    .then((snippets) => {
      dispatch({
        type: 'snippets/SNIPPETS_LOADED',
        snippets: snippets
      })
    })
  }
}
