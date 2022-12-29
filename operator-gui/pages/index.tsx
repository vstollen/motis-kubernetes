import { useState } from "react";
import FormCheckboxInput from "../components/FormCheckboxInput";
import FormCodeInput from "../components/FormCodeInput";
import FormTextInput from "../components/FormTextInput";

export default function Home() {
  const [periodicRefreshs, setPeriodicRefreshs] = useState(false);

  return (
    <div className="flex justify-center p-20">
      <section className="flex-auto max-w-2xl bg-zinc-50 rounded-md p-16">
        <div className="">
          <h1 className="mb-16 text-xl font-bold">Create MOTIS instance</h1>
          <div className="flex flex-col gap-8">
            <FormTextInput
              type="text"
              id="name"
              label="Name:"
              placeholder="motis-hessen"
            />
            <FormTextInput
              type="url"
              id="schedule"
              label="Schedule URL:"
              placeholder="https://download.geofabrik.de/europe/germany/nordrhein-westfalen/koeln-regbez-latest.osm.pbf"
            />
            <FormTextInput type="url" id="osm" label="Open Street Map (OSM):"
              placeholder="https://opendata.avv.de/current_GTFS/AVV_GTFS_mit_SPNV.zip"
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
            {periodicRefreshs && (
              <FormTextInput
                type="text"
                id="refreshCron"
                placeholder="0 3 * * *"
              />
            )}
          </div>
          <hr className="my-8" />
          <FormCodeInput id="config" label="config.ini:" />
          <button className="float-right mt-8 rounded-md bg-rose-600 hover:bg-rose-500 p-4 py-3 font-bold text-[hsl(347_100%_99%)] shadow shadow-rose-600/50 hover:shadow-rose-500/50 transition-all hover:shadow-md active:shadow">
            Create MOTIS instance
          </button>
        </div>
      </section>
    </div>
  );
}
