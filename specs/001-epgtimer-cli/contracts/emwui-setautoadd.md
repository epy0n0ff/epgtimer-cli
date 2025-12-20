# EMWUI API Contract: SetAutoAdd

**Feature**: `001-epgtimer-cli`
**Created**: 2025-12-20
**API Endpoint**: SetAutoAdd - Create automatic recording rules

---

## Endpoint Specification

**URL**: `POST http://<host>:<port>/api/SetAutoAdd?id=0`
**Method**: POST
**Content-Type**: `application/x-www-form-urlencoded; charset=UTF-8`
**Authentication**: None (trusted network)

---

## Request Parameters

### User-Configurable Parameters

| Parameter | Required | Type | Description | Example |
|-----------|----------|------|-------------|---------|
| andKey | Yes | string | Search keywords (program title must contain) | "わたしが恋人になれるわけ" |
| notKey | No | string | Exclusion keywords (program title must not contain) | "推しエンタ" |
| serviceList | Yes | string[] | Channel list in "ONID-TSID-SID" format (multiple values) | "32736-32736-1024" |

### Default Parameters (from curl sample)

All other parameters use the following defaults:

```
addchg=1
titleOnlyFlag=1
dayList=on
startTime=00:00
endTime=01:00
dateList=
freeCAFlag=0
chkDurationMin=0
chkDurationMax=0
chkRecDay=6
presetID=65535
onid=
tsid=
sid=
eid=
ctok=98357b8eedf096855c1cb636303ab2af
recMode=1
tuijyuuFlag=1
priority=2
useDefMarginFlag=1
serviceMode=1
tunerID=0
suspendMode=0
batFilePath=
batFileTag=
```

---

## Example Request

Based on the provided curl sample:

```http
POST /api/SetAutoAdd?id=0 HTTP/1.1
Host: 192.168.1.11:5510
Content-Type: application/x-www-form-urlencoded; charset=UTF-8

addchg=1&andKey=%E3%82%8F%E3%81%9F%E3%81%97%E3%81%8C%E6%81%8B%E4%BA%BA%E3%81%AB%E3%81%AA%E3%82%8C%E3%82%8B%E3%82%8F%E3%81%91&notKey=%E6%8E%A8%E3%81%97%E3%82%A8%E3%83%B3%E3%82%BF&titleOnlyFlag=1&serviceList=32736-32736-1024&serviceList=32736-32736-1025&serviceList=32737-32737-1032&dayList=on&startTime=00%3A00&endTime=01%3A00&dateList=&freeCAFlag=0&chkDurationMin=0&chkDurationMax=0&chkRecDay=6&presetID=65535&onid=&tsid=&sid=&eid=&ctok=98357b8eedf096855c1cb636303ab2af&recMode=1&tuijyuuFlag=1&priority=2&useDefMarginFlag=1&serviceMode=1&tunerID=0&suspendMode=0&batFilePath=&batFileTag=
```

**Key points**:
- Japanese text is URL-encoded (UTF-8)
- `serviceList` can appear multiple times for multiple channels
- Empty parameters (like `dateList=`) are included

---

## Response Format

**Success Response**:
- HTTP Status: 200 OK
- Body: (response format to be verified - may be JSON or plain text)

**Error Response**:
- HTTP Status: 400/500 or 200 with error message
- Detection: Parse response body for error indicators

---

## Client Implementation Notes

1. **Form Data Encoding**: Use `url.Values` in Go to build form data
2. **Multiple serviceList**: Call `v.Add("serviceList", value)` for each channel
3. **Japanese Text**: Automatic URL encoding when using `url.Values.Encode()`
4. **Empty Values**: Include empty string values as shown in curl sample

---

**Reference**: See data-model.md for Go struct definitions and ToFormData() method

Last Updated: 2025-12-20
