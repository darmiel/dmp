import { useState } from "react"
import DatePicker from "react-datepicker"

import "react-datepicker/dist/react-datepicker.css"

import { BsTriangle } from "react-icons/bs"
import { toast } from "sonner"

import { extractErrorMessage, PickerCustomInput } from "@/api/util"
import Button from "@/components/ui/Button"
import ModalContainer from "@/components/ui/modal/ModalContainer"
import { useAuth } from "@/contexts/AuthContext"

export default function CreateMeeting({
  projectID,
  onClose,
}: {
  projectID: number
  onClose: (newMeetingID: number) => void
}) {
  const [meetingTitle, setMeetingTitle] = useState<string>("")
  const [meetingDescription, setMeetingDescription] = useState<string>("")
  const [meetingDate, setMeetingDate] = useState<Date>(new Date())

  const { meetings: meeting } = useAuth()

  const createMeetingMutation = meeting!.useCreate(
    projectID,
    ({ data }, { __should_close }) => {
      toast.success(`Meeting #${data.ID} Created`)

      // clear form
      setMeetingTitle("")
      setMeetingDescription("")
      setMeetingDate(new Date())

      __should_close && onClose?.(data.ID)
    },
  )

  function create(shouldClose: boolean) {
    if (createMeetingMutation.isLoading) {
      return
    }
    createMeetingMutation.mutate({
      title: meetingTitle,
      description: meetingDescription,
      date: meetingDate,
      __should_close: shouldClose,
    })
  }

  return (
    <ModalContainer title="Create Meeting">
      <div className="space-y-2">
        <label className="text-neutral-400" htmlFor="meetingTitle">
          Meeting Title
        </label>
        <input
          id="meetingTitle"
          type="text"
          className="w-full rounded-lg border border-neutral-600 bg-neutral-800 p-2"
          placeholder="My awesome Meeting"
          onChange={(event) => setMeetingTitle(event.target.value)}
          value={meetingTitle}
          onKeyDown={(e) => e.key === "Enter" && create(false)}
        />
      </div>

      <div className="space-y-2">
        <label className="text-neutral-400" htmlFor="meetingDescription">
          Meeting Description
        </label>
        <textarea
          id="meetingDescription"
          className="w-full rounded-lg border border-neutral-600 bg-neutral-800 p-2"
          placeholder="(Markdown is supported)"
          rows={10}
          onChange={(event) => setMeetingDescription(event.target.value)}
          value={meetingDescription}
        />
      </div>

      <div className="space-y-2">
        <label className="text-neutral-400" htmlFor="topicDate">
          Topic Date
        </label>
        <div>
          <DatePicker
            selected={meetingDate}
            onChange={(date) => date && setMeetingDate(date)}
            timeInputLabel="Time:"
            dateFormat="MM/dd/yyyy h:mm aa"
            customInput={<PickerCustomInput />}
            showTimeInput
          />
        </div>
      </div>

      <hr className="border-neutral-600" />

      {createMeetingMutation.isError && (
        <div className="flex items-center space-x-2 text-sm font-bold text-red-500">
          <div>
            <BsTriangle />
          </div>
          <div>{extractErrorMessage(createMeetingMutation.error)}</div>
        </div>
      )}
      <div className="flex flex-row justify-end space-x-4">
        <Button
          style="secondary"
          isLoading={createMeetingMutation.isLoading}
          onClick={() => create(true)}
        >
          Save and Close
        </Button>

        <Button
          style="primary"
          isLoading={createMeetingMutation.isLoading}
          onClick={() => create(false)}
        >
          Save and Create New
        </Button>
      </div>
    </ModalContainer>
  )
}
