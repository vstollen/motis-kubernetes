import CardLayout from "../components/layouts/CardLayout";
import useSWR from "swr";
import { ListMotisInstance } from "../types/api";
import Button from "../components/Button";
import { useRouter } from "next/router";
import Link from "next/link";

const fetcher = (url: string) => fetch(url).then((res) => res.json());

export default function ListPage() {
  const { instances } = useMotisInstances();
  const router = useRouter();

  const instanceElements = instances
    ? instances.map((instance: ListMotisInstance) => {
        return (
          <li>
            <strong>{instance.name}</strong> – status: {instance.status}
          </li>
        );
      })
    : [];

  return (
    <CardLayout>
      <>
        <h1 className="mb-16 text-xl font-bold">Your MOTIS instances</h1>
        <ul className="list-inside list-disc">{instanceElements}</ul>
        <Button type="primary" onClick={() => router.push("/create")}>
          Create Instance
        </Button>
      </>
    </CardLayout>
  );
}

function useMotisInstances() {
  const { data } = useSWR("/api/motis", fetcher);

  return {
    instances: data ? data.instances : undefined,
  };
}
