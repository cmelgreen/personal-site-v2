// const [vis, reveal] = useReveal();
  //  const visited = useVisited();
  // const timeout = visited ? 0 : 3000;
  //  const wasScrolled = useScrolled();
  // className={classNames({
  //   'classes.animateReveal': visited || (!visited && wasScrolled)
  // })} 
  
// REVEAL
        // {/* <CardMedia 
      //   height="400"
      //   image={image}
      // /> */}
      // {/* <img
      //   className={classNames("big-picture", {
      //     "animate-reveal": visited || (!visited && wasScrolled)
      //   })}
      //   ref={ref}
      //   src={require("./big_picture.jpg")}
      //   alt="icon"
      // /> */}
      // {/* <div className="center-content big-picture-text" style={{ height: 200 }}>

      // </div> */}

      // const useScrolled = () => {
//   const trigger = useScrollTrigger({ threshold: 0.1 });
//   const [scrolled, setScrolled] = useState(false);

//   useEffect(() => {
//     if (trigger) {
//       setScrolled(true);
//     }
//   }, [trigger]);

//   return scrolled;
// };

// const useVisited = () => {
//   const [value, setValue] = useState(false);

//   useEffect(() => {
//     setValue(getCookie("visited"));
//     setCookie("visited", true);
//   }, []);

//   return value;
// };

// const getCookie = (cookie) => {
//   let value = document.cookie.match("(^|;)\\s*" + cookie + "\\s*=\\s*([^;]+)");
//   return value ? value.pop() : "";
// };

// const setCookie = (cookie, value) => {
//   document.cookie = `${cookie}=${value}`;
// };

// const useReveal = () => {
//     const [visible, setVisible] = useState('hidden')
//     const reveal = () => {
//       console.log("revealing")
//       setVisible("visible")
//     }
  
//     return [visible, reveal]
//   }