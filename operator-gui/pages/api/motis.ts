// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from "next";
import { customObjectsApi } from "../../services/kubernetes/kubernetes";

type Data = {
  message: string
}

type createMotisRequestBody = {
  name: string
  scheduleUrl: string
  osmUrl: string
  refreshSchedule: string
  config: string
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
  handlePostRequest(req, res);
}

async function handlePostRequest(req: NextApiRequest, res: NextApiResponse<Data>) {
  const requestBody: createMotisRequestBody = req.body;
  console.log(requestBody);
  await customObjectsApi.createNamespacedCustomObject("motis.motis-project.de", "v1alpha1", "default", "motis", {
    apiVersion: "motis.motis-project.de/v1alpha1",
    kind: "Motis",
    metadata: {
      name: requestBody.name
    },
    spec: {
      config: requestBody.name,
      items: [
        {
          key: "config-file",
          path: "config.ini"
        },
        {
          key: "schedules",
          path: "schedules"
        },
        {
          key: "osm",
          path: "osm"
        },
      ]
    }
  });
  res.status(202).json({ message: "success" });
}
