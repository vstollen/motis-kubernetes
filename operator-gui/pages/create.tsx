import { useState } from "react";
import FormCheckboxInput from "../components/FormCheckboxInput";
import FormCodeInput from "../components/FormCodeInput";
import FormTextInput from "../components/FormTextInput";
import CardLayout from "../components/layouts/CardLayout";
import Button from "../components/Button";
import {router} from "next/client";

export default function CreatePage() {
  const [name, setName] = useState("");
  const [scheduleUrl, setScheduleUrl] = useState("");
  const [osmUrl, setOsmUrl] = useState("");
  const [periodicRefreshs, setPeriodicRefreshs] = useState(false);
  const [refreshCron, setRefreshCron] = useState("");
  const [configIni, setConfigIni] = useState("");

  function handleSubmit() {
    fetch("/api/motis", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        scheduleUrl: scheduleUrl,
        osmUrl: osmUrl,
        refreshSchedule: periodicRefreshs ? refreshCron : null,
        config: configIni,
      }),
    });
  }

  return (
    <CardLayout>
      <>
        <h1 className="mb-16 text-xl font-bold">Create MOTIS instance</h1>
        <div className="flex flex-col gap-8">
          <FormTextInput
            type="text"
            id="name"
            value={name}
            onChange={setName}
            label="Name:"
            placeholder="motis-hessen"
          />
          <FormTextInput
            type="url"
            id="schedule"
            value={scheduleUrl}
            onChange={setScheduleUrl}
            label="Schedule URL:"
            placeholder="https://opendata.avv.de/current_GTFS/AVV_GTFS_mit_SPNV.zip"
          />
          <FormTextInput
            type="url"
            id="osm"
            value={osmUrl}
            onChange={setOsmUrl}
            label="Open Street Map (OSM):"
            placeholder="https://download.geofabrik.de/europe/germany/nordrhein-westfalen/koeln-regbez-latest.osm.pbf"
          />
        </div>
        <hr className="my-8" />
        <div className="flex flex-col">
          <FormCheckboxInput
            id="periodicRefreshs"
            label="Periodic Refreshs"
            checked={periodicRefreshs}
            onChange={() => setPeriodicRefreshs(!periodicRefreshs)}
          />
          <div className={periodicRefreshs ? "" : "hidden"}>
            <FormTextInput
              type="text"
              id="refreshCron"
              value={refreshCron}
              onChange={setRefreshCron}
              placeholder="0 3 * * *"
            />
          </div>
        </div>
        <hr className="my-8" />
        <FormCodeInput
          id="config"
          label="config.ini:"
          value={configIni}
          onChange={setConfigIni}
        />
        <Button type={"primary"} onClick={() => handleSubmit()}>Create MOTIS instance</Button>
        <Button type={"tertiary"} onClick={() => router.push("/")}>Cancel</Button>
      </>
    </CardLayout>
  );
}
