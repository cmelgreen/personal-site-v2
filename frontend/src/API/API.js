import axios from 'axios';

const apiRoot = "http://api.cmelgreen.com/"
const apiPostSummaries = apiRoot + "post-summaries"
const apiPost = apiRoot + "post/"

export const usePostSummaries = (numPosts=10) => {
  const [posts, setPosts] = useState([])

  useEffect(() => {
    axios.get(apiPostSummaries, {params: {numPosts}})
      .then(resp => {
        console.log(resp)
        if ( resp.data.posts ) setPosts(resp.data.posts)})
      .catch(() => setPosts([]))
  }, [])

  return posts
}

const api = axios.create({
  baseURL: apiPost,
  validateStatus: (status) => {
      return status == 200;
  }
});

const createPost = (post, idToken) => {
  return api.post(apiPost, post, {
    headers: {
      'Authorization': `Bearer ${idToken}` 
  }})
}

const updatePost = (post, idToken) => {
  return api.put(apiPost, post, {
    headers: {
      'Authorization': `Bearer ${idToken}` 
  }})
}

const usePostBySlug = (slug, raw=false) => {
  const [post, setPost] = useState({})

  useEffect(() => {
    if (slug) {
      api.get(apiPost + slug, {params: {raw}})
        .then(resp => setPost(resp.data))
        .catch(resp => setPost({}))
    }
  }, [slug])

  return [post, setPost]
}

const deletePost = (post, idToken) => {
  return api.delete(apiPost + post.slug, {
    headers: {
      'Authorization': `Bearer ${idToken}` 
  }})
}