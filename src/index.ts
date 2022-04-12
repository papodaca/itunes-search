import express from 'express';
import path from 'path';
import { fileURLToPath } from 'url';

import fetch from 'node-fetch';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const apple_base = "https://itunes.apple.com/search?entity=podcast&term=";

type FeedItem = {
  title: string,
  url: string
}

type ItunesItem = {
  collectionName: string,
  feedUrl: string
}

type ItunesResponce = {
  results: ItunesItem[]
}

async function run() {
  const app = express();
  app.set('views', path.join(__dirname, 'views'));
  app.set('view engine', 'hbs');

  app.get("/", async (req, res) => {
    let items: FeedItem[] = [];
    if(typeof req.query["query"] === "string") {
      let full_url = `${apple_base}${encodeURI(req.query["query"])}`
      let res2 = await fetch(full_url);
      let data = await res2.json() as ItunesResponce;
      items = data.results.map(item => ({
        title: item.collectionName,
        url: item.feedUrl
      }))
    }
    res.render('index', {
      query: req.query["query"],
      items
    })
  });

  const port = process.env.PORT || 3000
  app.listen(port, () => {
    console.log(`Listening on ${port}`);
  });
}

run().then(console.log, console.error);
