// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'

type Data = {
  message: string
}

type ApiError = {
  error: string
}

export default function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data | ApiError>
) {
  if (req.method != "POST") {
    res.status(405).json({ error: "HTTP method not allowed"});
    return;
  }
  res.status(200).json({ message: 'success' })
}
