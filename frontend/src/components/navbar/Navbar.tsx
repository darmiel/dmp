import Image from "next/image"
import Link from "next/link"
import { useRouter } from "next/router"
import { BsGithub } from "react-icons/bs"
import { ClipLoader } from "react-spinners"

import PerplexLogo from "@/../public/perplex.svg"
import { extractErrorMessage } from "@/api/util"
import UserAvatar from "@/components/user/UserAvatar"
import { useAuth } from "@/contexts/AuthContext"

export default function Navbar() {
  const { user, projects } = useAuth()

  const router = useRouter()
  const { project_id: projectID } = router.query

  const projectListQuery = projects!.useList()

  return (
    <aside
      id="separator-sidebar"
      className="flex-none top-0 left-0 w-20 h-screen transition-transform -translate-x-full sm:translate-x-0"
      aria-label="Sidebar"
    >
      <div className="h-full px-3 py-4 overflow-y-auto bg-black flex flex-col justify-between">
        <div>
          <div className="flex justify-center w-full mt-4">
            <Link href="/">
              <Image src={PerplexLogo} alt="Perplex Logo" />
            </Link>
          </div>

          <ul className="space-y-4 font-medium mt-10 flex flex-col items-center">
            {projectListQuery.isLoading ? (
              <ClipLoader color="white" />
            ) : projectListQuery.isError ? (
              <div>
                Error: <pre>{extractErrorMessage(projectListQuery.error)}</pre>
              </div>
            ) : (
              projectListQuery.data.data.map((project, key) => (
                <Link key={key} href={`/project/${project.ID}`}>
                  {String(project.ID) === projectID ? (
                    <li className="border-2 border-primary-900 rounded-full p-2">
                      <UserAvatar userID={String(project.ID)} />
                    </li>
                  ) : (
                    <li>
                      <UserAvatar userID={String(project.ID)} />
                    </li>
                  )}
                </Link>
              ))
            )}
          </ul>
        </div>

        <div className="mt-auto">
          <ul className="space-y-4 font-medium mt-10 flex flex-col items-center">
            <li className="p-2">
              <a href="https://github.com/darmiel/perplex">
                <BsGithub size="100%" />
              </a>
            </li>
            <li className="w-full">
              <hr className="border-neutral-600" />
            </li>
            <li>
              <Link href="/user">
                <UserAvatar userID={user?.uid ?? "1"} />
              </Link>
            </li>
          </ul>
        </div>
      </div>
    </aside>
  )
}
