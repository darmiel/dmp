import { useState } from "react"

import { extractErrorMessage } from "@/api/util"
import { useAuth } from "@/contexts/AuthContext"

import "reactjs-popup/dist/index.css"

import {
  Accordion,
  AccordionItem,
  Button,
  Checkbox,
  Chip,
  ScrollShadow,
} from "@nextui-org/react"
import clsx from "clsx"
import Link from "next/link"
import { useRouter } from "next/router"
import {
  BsArrowDown,
  BsArrowLeft,
  BsArrowUp,
  BsPlusCircle,
} from "react-icons/bs"
import { toast } from "sonner"

import { Topic } from "@/api/types"
import CreateTopicModal from "@/components/modals/TopicCreateModal"
import BadgeHeader from "@/components/ui/BadgeHeader"
import Flex from "@/components/ui/layout/Flex"
import ModalPopup from "@/components/ui/modal/ModalPopup"

function TopicListItem({
  topic,
  projectID,
  selected = false,
  beforeBeforeTopicID,
  beforeTopicID,
  afterAfterTopicID,
  afterTopicID,
}: {
  topic: Topic
  projectID: number
  selected?: boolean
  beforeBeforeTopicID?: number
  beforeTopicID?: number
  afterTopicID?: number
  afterAfterTopicID?: number
}) {
  const { user, topics } = useAuth()
  const toggleTopicMutation = topics!.useStatus(
    projectID,
    topic.meeting_id,
    () => {},
  )

  const updateOrderMutation = topics!.useUpdateOrder(
    projectID,
    topic.meeting_id,
    () => toast.success("Topic order updated"),
  )

  const isAssignedToUser = topic.assigned_users.some((u) => u.id === user?.uid)
  return (
    <Flex
      col
      className={clsx(
        "rounded-md border border-transparent p-2 transition duration-150 ease-in-out",
        "hover:border-neutral-800 hover:bg-neutral-900",
        {
          "border-l-3 border-l-primary-500": isAssignedToUser,
          "bg-neutral-800": selected,
        },
      )}
    >
      <div className="ml-1">
        <Flex justify="between">
          <Flex>
            <Checkbox
              isIndeterminate={toggleTopicMutation.isLoading}
              color={toggleTopicMutation.isLoading ? "warning" : "default"}
              isSelected={topic.closed_at.Valid}
              onValueChange={(checked) => {
                toggleTopicMutation.mutate({
                  topicID: topic.ID,
                  close: checked,
                })
              }}
              lineThrough={topic.closed_at.Valid}
            />
            <h2 className={clsx("truncate text-neutral-200", {})}>
              {topic.title}
              <span className="text-sm text-neutral-500"> #{topic.ID}</span>
            </h2>
          </Flex>
          <Flex>
            <Button
              isIconOnly
              startContent={<BsArrowUp />}
              size="sm"
              variant="light"
              onClick={(event) => {
                event.preventDefault()
                event.stopPropagation()
                updateOrderMutation.mutate({
                  topicID: topic.ID,
                  before: afterTopicID || -1,
                  after: afterAfterTopicID || -1,
                })
              }}
              isDisabled={
                afterAfterTopicID === undefined && afterTopicID === undefined
              }
            />
            <Button
              isIconOnly
              startContent={<BsArrowDown />}
              size="sm"
              variant="light"
              onClick={(event) => {
                event.preventDefault()
                event.stopPropagation()
                updateOrderMutation.mutate({
                  topicID: topic.ID,
                  before: beforeBeforeTopicID || -1,
                  after: beforeTopicID || -1,
                })
              }}
              isDisabled={
                beforeTopicID === undefined && beforeBeforeTopicID === undefined
              }
            />
          </Flex>
        </Flex>
        {topic.tags?.length > 0 && (
          <ScrollShadow orientation="horizontal" hideScrollBar className="mt-1">
            <Flex gap={1}>
              {topic.tags.map((tag) => (
                <Chip
                  key={tag.ID}
                  className="whitespace-nowrap"
                  variant="bordered"
                  style={{
                    borderColor: tag.color,
                  }}
                >
                  {tag.title}
                </Chip>
              ))}
            </Flex>
          </ScrollShadow>
        )}
      </div>
    </Flex>
  )
}

export default function TopicList({
  projectID,
  meetingID,
  selectedTopicID,
  onCollapse,
}: {
  projectID: number
  meetingID: number
  selectedTopicID?: number
  onCollapse?: () => void
}) {
  const [showCreateTopic, setShowCreateTopic] = useState(false)

  const router = useRouter()
  const { topics } = useAuth()

  const topicListQuery = topics!.useList(projectID, meetingID)

  if (topicListQuery.isLoading) {
    return <div>Loading...</div>
  }
  if (topicListQuery.isError) {
    return (
      <div>
        Error: <pre>{extractErrorMessage(topicListQuery.error)}</pre>
      </div>
    )
  }

  const checkedTopicCount = topicListQuery.data.data.filter(
    (topic) => topic.closed_at.Valid,
  ).length

  const checkedTopicRatio = checkedTopicCount / topicListQuery.data.data.length

  const showTopicListWithFilter = (filter: (topic: Topic) => boolean) => {
    const filtered = topicListQuery.data.data.filter(filter)
    return filtered.map((topic, index) => {
      const beforeBefore =
        filtered.length > index + 2 ? filtered[index + 2] : null
      const before = filtered.length > index + 1 ? filtered[index + 1] : null
      const after = index > 0 ? filtered[index - 1] : null
      const afterAfter = index > 1 ? filtered[index - 2] : null
      return (
        <div key={topic.ID}>
          <Link
            href={`/project/${projectID}/meeting/${topic.meeting_id}/topic/${topic.ID}`}
          >
            <TopicListItem
              projectID={projectID}
              topic={topic}
              selected={selectedTopicID === topic.ID}
              beforeTopicID={before?.ID}
              beforeBeforeTopicID={beforeBefore?.ID}
              afterTopicID={after?.ID}
              afterAfterTopicID={afterAfter?.ID}
            />
          </Link>
        </div>
      )
    })
  }
  const openTopics = showTopicListWithFilter((topic) => !topic.closed_at.Valid)
  const closedTopics = showTopicListWithFilter((topic) => topic.closed_at.Valid)

  return (
    <ul className="flex h-full flex-grow flex-col space-y-4 overflow-y-auto">
      <Flex justify="between">
        <Button
          onClick={() => setShowCreateTopic(true)}
          variant="light"
          startContent={<BsPlusCircle color="gray" size="1em" />}
          className="w-fit"
        >
          Create Topic
        </Button>
        {onCollapse && (
          <Button
            onClick={onCollapse}
            variant="light"
            isIconOnly
            startContent={<BsArrowLeft />}
          />
        )}
      </Flex>

      {/* ProgressBar */}
      <div className="px-4">
        <div className="h-2.5 w-full rounded-full bg-gray-700">
          <div
            className="h-2.5 rounded-full bg-primary-600"
            style={{
              width: `${checkedTopicRatio * 100}%`,
            }}
          ></div>
        </div>
        <div className="mt-2 text-center text-gray-500">
          <span className="text-white">{checkedTopicCount}</span> /{" "}
          {topicListQuery.data.data.length} topics done
        </div>
      </div>

      <hr className="mb-6 mt-4 border-gray-700" />

      <ScrollShadow hideScrollBar className="h-full p-2">
        <Accordion selectionMode="multiple" defaultSelectedKeys="all">
          <AccordionItem
            key="open-topics"
            title={
              <BadgeHeader title="Open Topics" badge={openTopics.length} />
            }
          >
            <div className="gap-4 space-y-1">{openTopics}</div>
          </AccordionItem>
          <AccordionItem
            key="closed-topics"
            title={
              <BadgeHeader title="Closed Topics" badge={closedTopics.length} />
            }
          >
            <div className="gap-4 space-y-1">{closedTopics}</div>
          </AccordionItem>
        </Accordion>
      </ScrollShadow>

      {/* Create Topic Popup */}
      <ModalPopup open={showCreateTopic} setOpen={setShowCreateTopic}>
        <CreateTopicModal
          projectID={projectID}
          meetingID={meetingID}
          onClose={(newTopicID?: number) => {
            setShowCreateTopic(false)
            newTopicID &&
              router.push(
                `/project/${projectID}/meeting/${meetingID}/topic/${newTopicID}`,
              )
          }}
        />
      </ModalPopup>
    </ul>
  )
}
