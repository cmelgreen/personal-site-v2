import React from "react";
import { useState } from "react";

import { makeStyles } from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import Grid from '@material-ui/core/Grid'
import CardActions from "@material-ui/core/CardActions";
import CardActionArea from "@material-ui/core/CardActionArea";
import CardContent from "@material-ui/core/CardContent";
import CardHeader from "@material-ui/core/CardHeader";
import CardMedia from "@material-ui/core/CardMedia";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";

export default function ContentCard(props) {
  const popOn = props.popOn ? props.popOn : 5;
  const popOff = props.popOff ? props.popOff : 0;
  const [pop, setPop] = useState(popOff);

  const accentOn = "secondary";
  const accentOff = "textPrimary";
  const [titleColor, setTitleColor] = useState(accentOff);

  const styles = {
    card: {
      display: "flex",
      alignItems: "flex-start",
      //height: 230,
      margin: 10,
      border: 10,
      marginLeft: "auto",
      marginRight: "auto",
      textAlign: "left"
    },
    header: {
      paddingBottom: 0,
      cursor: "pointer",
      transition: "all .1s ease-in-out",
    },
    title: {
      textAlign: 'left',
      variant: 'h5',

      
    },
    content: {
      paddingTop: 0,
      paddingBottom: 0,
      overflow: "hidden"
    },
    media: {
      margin: 0,
      height: 150,
      width: 150,
      transition: "all .1s ease-in-out"
    },
    actionArea: {
      "&:hover $focusHighlight": {
        opacity: 0
      },
      "&:hover $media": {
        transform: "scale(1.02)"
      }
    },
    actions: {
      marginTop: 0,
    },
    focusHighlight: {}
  };

  const classes = makeStyles(() => styles)();

  return (
    <>
      <Card className={classes.card} elevation={pop}>
        <MouseInOut function={setPop} in={popOn} out={popOff}>
          <CardActionArea
            classes={{
              root: classes.actionArea,
              focusHighlight: classes.focusHighlight
            }}
          >
            <MouseInOut function={setTitleColor} in={accentOn} out={accentOff}>
              <Grid container spacing={0} alignItems='center' >
                <Grid item xs={8} >
                  <CardHeader
                    className={classes.header}
                    title={props.post.title}
                    titleTypographyProps={{
                      color: titleColor,
                    }}
                    subheader={props.post.category}
                    subheaderTypographyProps={{
                      color: titleColor,
                      //variant: 'h6'
                    }}
                    onClick={props.handleClick}
                  />
                  <CardContent className={classes.content}>
  
                    <Typography color="textSecondary">
                      {props.post.summary}
                    </Typography>

                  </CardContent>
                  </Grid>
                <Grid item xs={4}>
                  <CardMedia className={classes.media} image={props.post.media} />
                </Grid>
              </Grid>
            </MouseInOut>
          </CardActionArea>
          <CardActions className={classes.actions}>
            <Tags tags={props.post.tags} />
          </CardActions>
        </MouseInOut>
      </Card>
    </>
  );
}

const Tags = (props) => {
  if (props.tags) {
    return (
      <div className="tags">
        {props.tags.map((tag, i) => (
          <Button key={i}>{tag}</Button>
        ))}
      </div>
    );
  } else {
    return null;
  }
};

const MouseInOut = (props) => (
  <div
    {...props}
    onMouseEnter={() => props.function(props.in)}
    onMouseLeave={() => props.function(props.out)}
  />
);
