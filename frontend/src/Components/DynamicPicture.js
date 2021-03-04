import React from 'react'

export default function DynamicPicture(props) {
    const imgSrc = props.src.split('.')
    const imgExt = imgSrc.pop()
  
    const resize = size => `${imgSrc.join('.')}-${size}.${imgExt}`

    console.log(resize('xs'))
  
    return (
      <picture >
        {/* TO-DO: Change to a map or something */}
        <source media="(max-width: 200px)" srcSet={resize("xs")} />
        <source media="(max-width: 600px)" srcSet={resize("sm")} />
        <source media="(max-width: 960px)" srcSet={resize("md")} />
        <source media="(max-width: 1280px)" srcSet={resize("lg")}/>
        <source media="(max-width: 1920px)" srcSet={resize("xl")} />
        <img className={props.className} 
          src={props.img} 
          // srcSet={`
          //   ${resize("xs")} 200px
          //   ${resize("sm")} 600px
          //   ${resize("md")} 960px
          //   ${resize("lg")} 1280px x2
          //   ${resize("xl")} 19200px x3
          // `}
          />
      </picture>
    )
  }