import React from 'react'

//import { Link } from "react-router-dom";

import {
  Card,
  CardActions,
  CardContent,
  CardHeader,
  CardActionArea,
  Typography,
} from "@material-ui/core";


export default function CMSSummaryCard(props) {
  return (
    <div className="cms-summary-card">
      <Card>
      {/*component={Link} */}
        <CardActionArea className={"cardActionArea"}  >
            <CardHeader title={props.post.title} />
            <CardContent>
              <Typography color="textSecondary">
                {props.post.summary}
              </Typography>
            </CardContent>
        </CardActionArea>
      </Card>
    </div>
  )
}