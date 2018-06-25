# ETL Conventions for use with PEDSnet CDM V2.0 OMOP V5

The PEDSnet Common Data Model is an evolving specification, based in structure on the OMOP Common Data Model, but expanded to accommodate requirements of both the PCORnet Common Data Model and the primary research cohorts established in PEDSnet.

Version 2 of the PEDSnet CDM reflects the ETL processes developed after the first six months of network development. As such, it proposes to align with version 2 of the PCORnet CDM.

This document provides the ETL processing assumptions and conventions developed by the PEDSnet data partners that should be used by a data partner for ensuring common ETL business rules. This document will be modified as new situations are identified, incorrect business rules are identified and replaced, as new analytic use cases impose new/different ETL rules, and as the PEDSnet CDM continues to evolve.

Comments on this specification and ETL rules are welcome. Please send email to pedsnetdcc@email.chop.edu, or contact the PEDSnet project management office (details available via http://www.pedsnet.info).

#### PEDSnet Data Standards and Interoperability Policies:

1. The PEDSnet data network will store data using structures compatible with the PEDSnet Common Data Model (PCDM).

2. The PCDM v2 is based on the Observational Medical Outcomes Partnership (OMOP) data model, version 5. 

3. A subset of data elements in the PCDM will be identified as principal data elements (PDEs). The PDEs will be used for population-level queries. Data elements which are NOT PDEs will be marked as Optional (ETL at site discretion) or Non-PDE (ETL required, but data need not be transmitted to DCC), and will not be used in queries without prior approval of site.

4. It is anticipated that PEDSnet institutions will make a good faith attempt to obtain as many of the data elements not marked as Optional as possible.

5. The data elements classified as PDEs and those included in the PCDM will be approved by the PEDSnet Executive Committee (comprised of each PEDSnet institution's site principal investigator).

6. Concept IDs are taken from OMOP 5 vocabularies for PEDSnet CDM v2, using the complete (restricted) version that includes licensed terminologies such as CPT and others.

7. PCORnet CDM V2.0 requires data elements that are not currently considered "standard concepts". Vocabulary version 5 has a new vocabulary (vocabulary_id=PCORNet) that was added by OMOP to capture all of the PCORnet concepts that are not in the standard terminologies. We use concept_ids from vocabulary_id=PCORNet where there are no existing standard concepts. We highlight where we are pulling concept_ids from vocabulary_id=PCORNet in the tables. While terms from vocabulary_id=PCORNet violates the OMOP rule to use only concept_ids from standard vocabularies vocabulary_id=PCORNet is a non-standard vocabulary), this convention enables a clean extraction from PEDSnet CDM to PCORnet CDM.

8. Some source fields may be considered sensitive by data sites. Potential examples include patient_source_value, provider_source_value, care_site_source_value. Many of these fields are used to generate an ID field, such as PERSON.patient_source_value PERSON.person_id, that is used as a primary key in PERSON and a foreign key in many other tables. Sites are free to obfuscate or not provide source values that are used to create ID variables. Sites must maintain a mapping from the ID variable back to the original site-specific value for local re-identification tasks.

    1. Source fields that contain clinical data, such as source condition occurrence, should be included

    2. The PEDSnet DCC will never release source values to external data partners.

    3. Source value obfuscation techniques may include replacing the real source value with a random number, an encrypted derivative value/string, or some other site-specific algorithm.

9. Regarding the nullability of all source value (string) fields only, the PEDSnet CDM will accommodate the following values, taken from the PCORnet CDM:

Null Name | Definition of each field
 --- | ---
NULL | A data field is not present in the source system. Note. This is not a 'NULL' string but the NULL value.
'NI' = No Information | A data field is present in the source system, but the source value is null or blank
'UN' = Unknown | A data field is present in the source system, but the source value explicitly denotes an unknown value
'OT' = Other | A data field is present in the source system, but the source value cannot be mapped to the CDM

***ETL Recommendation:*** Due to PK/FK constraints, the most efficient order for ETL table is location, care_site, provider, person, visit_occurrence, condition_occurrence, observation, procedure_occurrence,measurement,drug exposure and observation_period

It is recommended to refer to the vocabulary documentation as provided by ODHSII for guidance on how to populate "concept_id" fields in the model and for any specific transformations in the vocabulary. http://www.ohdsi.org/web/wiki/doku.php?id=documentation:vocabulary:data_etl

* * *
## Table of Contents
####[1.1 Person](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#11-person-1)

####[1.2 Death](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#12-death-1)

####[1.3 Location](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#13-location-1)

####[1.4 Caresite](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#14-care_site)

####[1.5 Provider](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#15-provider-1)

####[1.6 Visit Occurrence ](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#16-visit_occurrence)

####[1.7 Condition Occurrence](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#17-condition_occurrence)

####[1.8 Procedure Occurrence](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#18-procedure_occurrence)

####[1.9 Observation](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#19-observation-1)

####[1.10 Observation Period](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#110-observation-period-1)

####[1.11 Drug Exposure](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#111-drug-exposure-draft)

####[1.12 Measurement](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#112-measurement-draft)

####[1.13 Fact Relationship](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#113-fact-relationship-1)

####[1.14 Visit Payer](Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#114-visit_payer)

####[Appendix] (Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#a1-abms-specialty-category-to-omop-v5-specialty-mapping)

* * *
## 1.1 PERSON

The person domain contains records that uniquely identify each patient in the source data who is time at-risk to have clinical observations recorded within the source systems. Each person record has associated demographic attributes, which are assumed to be constant for the patient throughout the course of their periods of observation. All other patient-related data domains have a foreign-key reference to the person domain.

PEDSnet uses a specific definition of an active PEDSnet patient. Only patients who meet the PEDSnet definition of an active patient should be included in this table. The criteria for identifying an active patient are:

- Has a unique identifier AND

- At least 1 "in person" clinical encounter on or after January 1, 2009 AND

- At least 1 coded diagnoses recorded on or after January 1, 2009 AND

- Is not a test patient or a research-only patient

The definition of an "in person" clinical encounter remains heuristic -any encounter type that involves a meaningful \*\*physical\*\* interaction with a clinician that involved clinical content. An encounter for a suture removal or a telephone encounter or a lab blood draw does not meet this definition.

**NOTE: While the 1/1/2009 date and "in person" clinical encounter restrictions apply to defining an active PEDSnet patient, once a patient has met this criteria, PEDSnet will extract *ALL* available clinical encounters/clinical data of any type across all available dates. That is, 1/1/2009 and 1 'in person' clinical encounter applies only to defining the active patient cohort. It does NOT apply to data extraction on active patients.**

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
person_id | Yes* | Integer | A unique identifier for each person; this is created by each contributing site. | <p>This is not a value found in the EHR.</p> PERSON_ID must be unique for all patients within a single data set.</p><p>Sites may choose to use a sequential value for this field
gender_concept_id | Yes | Integer |  A foreign key that refers to a standard concept identifier in the Vocabulary for the gender of the person. | Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table select \* from concept where domain_id='Gender'): <ul><li>Ambiguous: concept_id = 44814664 </li> <li>Female: concept_id = 8532</li> <li>Male: concept_id = 8507</li> <li>No Information: concept_id = 44814650 (Vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653</li> <li>Other: concept_id = 44814649</li></ul>
gender_source_concept_id | No | Integer | A foreign key to the gender concept that refers to the code used in the source.
year_of_birth | Yes* | Integer |  The year of birth of the person. | <p>For data sources with date of birth, the year is extracted. For data sources where the year of birth is not available, the approximate year of birth is derived based on any age group categorization available.</p> Please keep all accurate/real dates (No date shifting)
month_of_birth | No* | Integer | The month of birth of the person. | <p>For data sources that provide the precise date of birth, the month is extracted and stored in this field.</p> Please keep all accurate/real dates (No date shifting)
day_of_birth | No | Integer | The day of the month of birth of the person. | <p>For data sources that provide the precise date of birth, the day is extracted and stored in this field.</p> Please keep all accurate/real dates (No date shifting)
time_of_birth | No* | Datetime |  The birth date and time | <p>Do not include timezone.</p> Please keep all accurate/real dates (No date shifting). If there is no time associated with the date assert midnight.
race_concept_id | No* | Integer |  A foreign key that refers to a standard concept identifier in the Vocabulary for the race of the person. | Details of categorical definitions: <ul><li>**-American Indian or Alaska Native**: A person having origins in any of the original peoples of North and South America (including Central America), and who maintains tribal affiliation or community attachment.</li> <li>**-Asian**: A person having origins in any of the original peoples of the Far East, Southeast Asia, or the Indian subcontinent including, for example, Cambodia, China, India, Japan, Korea, Malaysia, Pakistan, the Philippine Islands, Thailand, and Vietnam.</li> <li>**-Black or African American**: A person having origins in any of the black racial groups of Africa.</li> <li>**-Native Hawaiian or Other Pacific Islander**: A person having origins in any of the original peoples of Hawaii, Guam, Samoa, or other Pacific Islands.</li> <li>**-White**: A person having origins in any of the original peoples of Europe, the Middle East, or North Africa.</li></ul> <p>For patients with multiple races (i.e. biracial), race is considered a single concept, meaning there is only one race slot. If there are multiple races in the source system, concatenate all races into one race_source_value (see below) and use concept_id code as 'Multiple Race.'</p> Predefined values (valid concept_ids found in CONCEPT table where (domain_id='Race' and vocabulary_id = 'Race') or (vocabulary_id = 'PCORNet' and concept_class_id='Race'): <ul><li>American Indian/Alaska Native: concept_id = 8657</li> <li>Asian: concept_id = 8515</li> <li>Black or African American: concept_id = 8516</li> <li>Native Hawaiian or Other Pacific Islander: concept_id = 8557</li> <li>White: concept_id = 8527</li> <li>Multiple Race: concept_id = 44814659 (vocabulary_id='PCORNet')</li> <li>Refuse to answer: concept_id = 44814660 (vocabulary_id='PCORNet')</li> <li>No Information: concept_id = 44814650 vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653</li> <li>Other: concept_id = 44814649</li></ul>
race_source_concept_id| No | Integer| A foreign key to the race concept that refers to the code used in the source.
ethnicity_concept_id | No | Integer | A foreign key that refers to the standard concept identifier in the Vocabulary for the ethnicity of the person. | <p>For PEDSnet, a person with Hispanic ethnicity is defined as "A person of Cuban, Mexican, Puerto Rican, South or Central American, or other Spanish culture or origin, regardless of race."</p> Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id ='Ethnicity' or vocabulary_id=PCORNet where noted): <ul><li>Hispanic: concept_id = 38003563</li> <li>Not Hispanic: concept_id = 38003564</li> <li>No Information: concept_id = 44814650 (vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653 (vocabulary_id='PCORNet')</li> <li>Other: concept_id = 44814649 (vocabulary_id='PCORNet')</li></ul>
ethnicity_source_concept_id | No | Integer | A foreign key to the ethnicity concept that refers to the code used in the source.
location_id | No |Integer |  A foreign key to the place of residency (ZIP code) for the person in the location table, where the detailed address information is stored.
provider_id | No |Integer |  Foreign key to the primary care provider the person is seeing in the provider table.| <p>For PEDSnet CDM V2.0: Sites will use site-specific logic to determine the best primary care provider and document how that decision was made (e.g., billing provider).</p>
care_site_id | Yes |Integer |  A foreign key to the site of primary care in the care_site table, where the details of the care site are stored | For patients who receive care at multiple care sites, use site-specific logic to select a care site that best represents where the patient obtains the majority of their recent care. If a specific site within the institution cannot be identified, use a care_site_id representing the institution as a whole.
pn_gestational_age | No |Integer |  The post-menstrual age in weeks of the person at birth, if known | Use granularity of age in weeks as is recorded in local EHR.
person_source_value | Yes | Varchar |  An encrypted key derived from the person identifier in the source data. | <p>Insert a unique pseudo-identifier (random number, encrypted identifier) into the field. Do not insert the actual MRN or PAT_ID from your site. A mapping from the pseudo-identifier for person_source_value in this field to a real patient ID or MRN from the source EHR must be kept at the local site. This mapping is not shared with the data coordinating center. It is used only by the site for re-identification for study recruitment or for data quality review.</p>
gender_source_value | No* | Varchar |  The source code for the gender of the person as it appears in the source data. | The person's gender is mapped to a standard gender concept in the Vocabulary; the original value is stored here for reference. See gender_concept_id
race_source_value | No* | Varchar |  The source code for the race of the person as it appears in the source data. | <p>The person race is mapped to a standard race concept in the Vocabulary and the original value is stored here for reference.</p> For patients with multiple races (i.e. biracial), race is considered a single concept, meaning there is only one race slot. If there are multiple races in the source system, concatenate all races into one source value, and use the concept_id for Multiple Race.
ethnicity_source_value | No* | Varchar |  The source code for the ethnicity of the person as it appears in the source data. | The person ethnicity is mapped to a standard ethnicity concept in the Vocabulary and the original code is, stored here for reference.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

## 1.2 DEATH

The death domain contains the clinical event for how and when a person dies. Living patients should not contain any information in the death table.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
person_id | Yes | Integer | A foreign key identifier to the deceased person. The demographic details of that person are stored in the person table.| See PERSON.person_id (primary key)
death_date | Yes | Date | The date the person was deceased. | <p>If the precise date including day or month is not known or not allowed, December is used as the default month, and the last day of the month the default day. If no date available, use date recorded as deceased.</p> When the date of death is not present in the source data, use the date the source record was created.
death_time | Yes | Datetime | The date the person was deceased. |<p>**This field is custom to PEDSnet**</p> <p>If the precise date including day or month is not known or not allowed, December is used as the default month, and the last day of the month the default day. If no date available, use date recorded as deceased.</p> <p>When the date of death is not present in the source data, use the date the source record was created. If there is no time associated with the date assert midnight.</p>
death_type_concept_id | Yes | Integer | A foreign key referring to the predefined concept identifier in the Vocabulary reflecting how the death was represented in the source data. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where domain_id ='Death Type')</p> <p>select \* from concept where domain_id ='Death Type' yields 9 valid concept_ids. If none are correct, use concept_id = 0</p> <p>Note: Most current ETLs are extracting data from EHR so most likely concept_id to insert here is 38003569 ("EHR record patient status "Deceased")</p> Note: These terms only describe the source from which the death was reported. It does not describe our certainty/source of the date of death, which may have been created by one of the heuristics described in death_date.
cause_concept_id | No | Integer | A foreign referring to a standard concept identifier in the Vocabulary for conditions. | 
cause_source_value | No | Varchar | The source code for the cause of death as it appears in the source. This code is mapped to a standard concept in the Vocabulary and the original code is stored here for reference.
cause_source_concept_id | No | Integer | A foreign key to the vocbaulary concept that refers to the code used in the source.| This links to the concept id of the vocabulary of the cause of death concept id as stored in the source. For example, if the cause of death is "Acute myeloid leukemia, without mention of having achieved remission" which has an icd9 code of 205.00 the cause source concept id is 44826430 which is the icd9 code concept that corresponds to the diagnosis 205.00.

#### 1.2.1 Additional Notes

- Each Person may have more than one record of death in the source data. It is OK to insert multiple death records for an individual.
- If the Death Date cannot be precisely determined from the data, the best approximation should be used.

## 1.3 LOCATION

The Location domain represents a generic way to capture physical location or address information. Locations are used to define the addresses for Persons and Care Sites. The most important field is ZIP for location-based queries.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
location_id | Yes* | Integer | A unique identifier for each geographic location. | This is not a value found in the EHR. Sites may choose to use a sequential value for this field
state | No | Varchar | The state field as it appears in the source data.
zip | No* | Varchar | The zip code. For US addresses, valid zip codes can be 3, 5 or 9 digits long, depending on the source data. | While optional, this is the most important field in this table to support location-based queries.
location_source_value | No | Varchar | <p>Optional - Do not transmit to DCC.</p> The verbatim information that is used to uniquely identify the location as it appears in the source data. | <p>If location source values are deemed sensitive by your organization, insert a pseudo-identifier (random number, encrypted identifier) into the field. Sites electing to obfuscate location_source_values will keep the mapping between the value in this field and the original clear text location source value. This value is only used for site-level re-identification for study recruitment and for data quality review.</p> Sites may consider using the location_id field value in this table as the pseudo-identifier as long as a local mapping from location_id to the real site identifier is maintained.
address_1 | No | Varchar | |Optional - Do not transmit to DCC
address_2 | No | Varchar | |Optional - Do not transmit to DCC
city | No | Varchar | |Optional - Do not transmit to DCC
county | No | Varchar | |Optional - Do not transmit to DCC

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.3.1 Additional Notes

- Each address or Location is unique and is present only once in the table
- Locations in this table are restricted to locations that are applicable to persons and care_sites in the Pedsnet cohort at each site. When external data is implemented, valid(data containing) locations may be expanded beyond locations of those only present in clinical tables.

## 1.4 CARE_SITE

The Care Site domain contains a list of uniquely identified physical or organizational units where healthcare delivery is practiced (offices, wards, hospitals, clinics, etc.).

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
care_site_id | Yes* | Integer | A unique identifier for each defined location of care within an organization. Here, an organization is defined as a collection of one or more care sites that share a single EHR database. | <p>This is not a value found in the EHR.</p> Sites may choose to use a sequential value for this field
care_site_name | No | Varchar | The description of the care site | 
place_of_service_concept_id | No |  Integer | A foreign key that refers to a place of service concept identifier in the Vocabulary | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where domain_id='Visit')</p> <p>select \* from concept where domain_id ='Visit' yields 4 valid concept_ids.</p> If none are correct, use concept_id = 0 <ul><li>Inpatient Hospital Stay: concept_id = 9201</li> <li>Ambulatory Visit: concept_id = 9202</li> <li>Emergency Department: concept_id = 9203</li> <li>Long Term Care Visit = 42898160 </li><li>Other ambulatory Visit = 44814711</li> <li>Non-Acute Institutional Stay: concept_id = 44814710 (vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653 (vocabulary_id='PCORNet')</li> <li>Other: concept_id =  44814649 (vocabulary_id='PCORNet')</li> <li>No information: concept_id =  44814650(vocabulary_id='PCORNet')</li></ul>
location_id | No* |  Integer | A foreign key to the geographic location of the administrative offices of the organization in the location table, where the detailed address information is stored.
care_site_source_value | Yes |  Varchar | The identifier for the organization in the source data, stored here for reference. | <p>If care site source values are deemed sensitive by your organization, insert a pseudo-identifier (random number, encrypted identifier) into the field. Sites electing to obfuscate care site_source_values will keep the mapping between the value in this field and the original clear text location source value. This value is only used for site-level re-identification for study recruitment and for data quality review.</p> <p>For EPIC EHRs, map care_site_id to Clarity Department.</p> Sites may consider using the care_site_id field value in this table as the pseudo-identifier as long as a local mapping from care_site_id to the real site identifier is maintained.
place_of_service_source_value | No | Varchar | The source code for the place of service as it appears in the source data, stored here for reference.
specialty_concept_id|No|Integer|The specialty of the department linked to a standard specialty concept as it appears in the Vocabulary | <p>Care sites could have one or more specialties or a Care site could have no specialty information.</p><p>**Valid specialty concept ids for PEDSnet are found in the [appendix] (https://github.com/PEDSnet/Data_Models/blob/master/PEDSnet/V2/docs/Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#a1-abms-specialty-category-to-omop-v5-specialty-mapping)**</p><p>**Please use the following rules:**</p><ul><li><p> If care site specialty information is unavaiable, please follow the convention on reporting values that are unknown,null or unavailable. </p></li><li><p> If a care site has a single specialty associated with it, sites should link the specialty to the **valid specialty concepts as assigned in the [appendix] (https://github.com/PEDSnet/Data_Models/blob/master/PEDSnet/V2/docs/Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#a1-abms-specialty-category-to-omop-v5-specialty-mapping)**. If the specialty does not correspond to a value in this listing, please use the NUCC Listing (vocabulary_id='NUCC') provided in the vocabulary as a reference. </p></li><li><p> If there are multiple specialties assicated with a particular care site and sites are not able to assign a specialty value on the visit occurrence level, sites should use the specialty concept id=38004477 "Pediatric Medicine". </p></li><li><p> If there are multiple specialties aossicated with a particular care site and this information is attainable, sites should document the strategy used to obtain this information and the strategy used to link the correct care site/specialty pair for each visit occurrence. Sites should also link the specialty to the **valid specialty concepts as assigned in the [appendix] (https://github.com/PEDSnet/Data_Models/blob/master/PEDSnet/V2/docs/Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#a1-abms-specialty-category-to-omop-v5-specialty-mapping)**</p> If the specialty does not correspond to a value in this listing, please use the NUCC Listing (vocabulary_id='NUCC') provided in the vocabulary as a reference. </li>
</ul>|
specialty_source_value| No | Varchar | The source code for the specialty as it appears in the source data, stored here for reference.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.4.1 Additional Notes

- Care sites are primarily identified based on the specialty or type of care provided, and secondarily on physical location, if available (e.g. North Satellite Endocrinology Clinic)
- The Place of Service Concepts are based on a catalog maintained by the CMS (see vocabulary for values)

## 1.5 PROVIDER

The Provider domain contains a list of uniquely identified health care providers. These are typically physicians, nurses, etc.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
provider_id | Yes |  Integer | A unique identifier for each provider. Each site must maintain a map from this value to the identifier used for the provider in the source data. | This is not a value found in the EHR. Sites may choose to use a sequential value for this field. See Additional Comments below. Sites should document who they have included as a provider.
provider_name | No | Varchar | A description of the provider | DO NOT TRASMIT TO DCC
gender_concept_id | No | Integer | The gender of the provider | A foreign key to the concept that refers to the code used in the source.|Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table select \* from concept where domain_id='Gender'): <ul><li>Ambiguous: concept_id = 44814664 </li> <li>Female: concept_id = 8532</li> <li>Male: concept_id = 8507</li> <li>No Information: concept_id = 44814650 (Vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653</li> <li>Other: concept_id = 44814649</li></ul>
specialty_concept_id | No |  Integer | A foreign key to a standard provider's specialty concept identifier in the Vocabulary. | <p>Please map the source data to the mapped provider specialtity concept associated with the American Medical Board of Specialties as seen in **Appendix A1**. Predefined value set (valid concept_ids found in CONCEPT table where domain_id='Provider Specialty' and  vocabulary_id = Specialty)</p> <p>select \* from concept where domain_id ='Provider Specialty' and vocabulary_id='Specialty' and invalid_reason is null yields 107 valid concept_ids.</p> <p>If none are correct, use concept_id = 0</p> For providers with more than one specialty, use site-specific logic to select one specialty and document the logic used. For example, sites may decide to always assert the \*\*first\*\* specialty listed in their data source. If the specialty does not correspond to a value in this listing, please use the NUCC Listing (vocabulary_id='NUCC') provided in the vocabulary as a reference.
care_site_id | Yes |  Integer | A foreign key to the main care site where the provider is practicing. | See CARE_SITE.care_site_id (primary key)
year_of_birth | No | Integer | The year of birth of the provider||
NPI | No |  Varchar | The National Provider Identifier (NPI) of the provider. | <p>Optional - Do not transmit to DCC.</p> 
DEA | No |  Varchar | The Drug Enforcement Administration (DEA) number of the provider. | <p>Optional - Do not transmit to DCC.</p>
provider_source_value | Yes |  Varchar | The identifier used for the provider in the source data, stored here for reference. | <p>Insert a pseudo-identifier (random number, encrypted identifier) into the field. Do not insert the actual PROVIDER_ID from your site. A mapping from the pseudo-identifier for provider_source_value in this field to a real provider ID from the source EHR must be kept at the local site. This mapping is not shared with the data coordinating center. It is used only by the site for re-identification for study recruitment or for data quality review.</p> Sites may consider using the provider_id field value in this table as the pseudo-identifier as long as a local mapping from provider_id to the real site identifier is maintained.
specialty_source_value | No |  Varchar | The source code for the provider specialty as it appears in the source data, stored here for reference. | Optional. May be obfuscated if deemed sensitive by local site.
specialty_source_concept_id | No | Integer | A foreign key to a concept that refers to the code used in the source.| If providing this information, sites should document how they determine the specialty associated with the provider. **Valid specialty concept ids for PEDSnet are found in the [appendix] (https://github.com/PEDSnet/Data_Models/blob/master/PEDSnet/V2/docs/Pedsnet_CDM_V2_OMOPV5_ETL_Conventions.md#a1-abms-specialty-category-to-omop-v5-specialty-mapping)** If the specialty does not correspond to a value in this listing, please use the NUCC Listing (vocabulary_id='NUCC') provided in the vocabulary as a reference
gender_source_value | No | Varchar | The source value for the provider gender.
gender_source_concept_id | No | Integer | The gender of the provider as represented in the source that maps to a concept in the vocabulary|  

#### 1.5.1 Additional Notes

- For PEDSnet, a provider is any individual (MD, DO, NP, PA, RN, etc) who is authorized to document care.
- Providers are not duplicated in the table.

## 1.6 VISIT_OCCURRENCE

The visit occurrence domain contains the spans of time a person continuously receives medical services from one or more providers at a care site in a given setting within the health care system.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
visit_occurrence_id | Yes* |  Integer | A unique identifier for each person's visits or encounter at a healthcare provider. | This is not a value found in the EHR. Sites may choose to use a sequential value for this field. Do not use institutional encounter ID.
person_id | Yes* |  Integer | A foreign key identifier to the person for whom the visit is recorded. The demographic details of that person are stored in the person table.
visit_start_date | Yes* |  Date | The start date of the visit. | No date shifting. Full date and time. If there is no time associated with the date assert midnight.
visit_end_date | No* | Date | The end date of the visit. | <p>No date shifting.</p> <p>If this is a one-day visit the end date should match the start date.</p> If the encounter is on-going at the time of ETL, this should be null. I Full date and time. If there is no time associated with the date assert midnight.
visit_start_time |Yes |  Datetime | The start date of the visit. | No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
visit_end_time | No | Datetime | The end date of the visit. | <p>No date shifting.</p> <p>If this is a one-day visit the end date should match the start date.</p> If the encounter is on-going at the time of ETL, this should be null.  Full date and time. If there is no time associated with the date assert midnight.
provider_id | No* |  Integer | A foreign key to the provider in the provider table who was associated with the visit. | <p>Use attending or billing provider for this field if available, even if multiple providers were involved in the visit. Otherwise, make site-specific decision on which provider to associate with visits and document.</p> **NOTE: this is NOT required in OMOP CDM v4, but appears in OMOP CDMv5.**
care_site_id | No* | Integer | A foreign key to the care site in the care site table that was visited. | See CARE_SITE.care_site_id (primary key)
visit_concept_id | Yes | Integer | A foreign key that refers to a place of service concept identifier in the vocabulary. | <p>**In PEDSnet CDM v1, this field was previously called place_of_service_concept_id**</p><p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id ='Visit' or (vocabulary_id='PCORNet' and concept_class_id='Encounter Type').</p> <p>select \* from concept where domain_id ='Visit' yields 4 valid concept_ids.</p> If none are correct, use concept_id = 0 <ul><li>Inpatient Hospital Stay: concept_id = 9201</li> <li>Ambulatory Visit: concept_id = 9202</li> <li>Emergency Department: concept_id = 9203</li> <li>Long Term Care Visit = 42898160 </li><li>Other ambulatory Visit = 44814711</li> <li>Non-Acute Institutional Stay: concept_id = 44814710 (vocabulary_id='PCORNet')</li> <li>Unknown: concept_id = 44814653 (vocabulary_id='PCORNet')</li> <li>Other: concept_id =  44814649 (vocabulary_id='PCORNet')</li> <li>No information: concept_id =  44814650(vocabulary_id='PCORNet')</li></ul>
visit_type_concept_id | Yes | Integer | A foreign key to the predefined concept identifier in the standard vocabulary reflecting the type of source data from which the visit record is derived.| <p> select \* from chop_omop5.concept where domain_id='Visit Type' yields 3 valid concept_ids.</p> If none are correct, user concept_id=0. The majoirty of visits should be type 'Visit derived from EHR record' which is concept_id=44818518
visit_source_value | No | Varchar | The source code used to reflect the type or source of the visit in the source data. Valid entries include office visits, hospital admissions, etc. These source codes can also be type-of service codes and activity type codes.
visit_source_concept_id | No* | Integer | A foreign key to a concept that refers to the code used in the source. | If a site is using HCPS or CPT for their visit source value, the standard concept id that maps to the particular vocabulary can be used here.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.6.1 Additional Notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is **NOT** applied to visit_occurrence. All visits, of all types (physical and virtual) are included for an active patient.
- A Visit Occurrence is recorded for each visit to a healthcare facility.
- **If a visit includes moving between different visit_concepts (ED -\> inpatient) this should be split into separate visit_occurrence records to meet PCORnet's definitions**.

**To show the relationship of the split (ED -\> inpatient) encounter, use the FACT_RELATIONSHIP table**.

An example of this is below:

**VISIT_OCCURRENCE**

visit_occurrence_id | person_id |    visit_start_date    |     visit_end_date     | provider_id | care_site_id | place_of_service_concept_id | place_of_service_source_value 
---------------------|-----------|---------|------------------------|-------------|--------------|------|---------------
35022489 |    209846 | 2011-11-14 17:36:00-05 | 2011-11-14 22:25:00-05 |        2238 |          322 |                        9203 | Emergency
35022490 |    209846 | 2011-11-14 22:25:00-05 | 2011-11-15 16:33:00-05 |        2238 |           43 |                        9201 | Emergency

**FACT_RELATIONSHIP**

Domain_concept_id_1 | fact_id_1 | Domain_concept_id_2 | fact_id_2 | relationship_concept_id
--- | --- | --- | --- | ---
Visit | 35022489| Visit | 35022490|  Occurs before
Visit | 35022490 | Visit | 35022489 | Occurs after 

Because the domain_concept_id and relationship_concept_id are actually numeric values the following is an example of how the table is stored:

Domain_concept_id_1 | fact_id_1 | Domain_concept_id_2 | fact_id_2 | relationship_concept_id
--- | --- | --- | --- | ---
8| 35022489| 8 | 35022490|  44818881
8 | 35022490 | 8 | 35022489 | 44818783

- Operating and Anesthesia encounters that occur as apart of the Inpatient stay should be rolled up into one Inpatient encounter.
- Each Visit is standardized by assigning a corresponding Concept Identifier based on the type of facility visited and the type of services rendered.
- At any one day, there could be more than one visit.
- One visit may involve multiple attending or billing providers (e.g. billing, attending, etc), in which case the ETL must specify how a single provider id is selected or leave the provider_id field null.
- One visit may involve multiple care sites, in which case the ETL must specify how a single care_site id is selected or leave the care_site_id field null.


## 1.7 CONDITION_OCCURRENCE

The condition occurrence domain captures records of a disease or a medical condition based on diagnoses, signs and/or symptoms observed by a provider or reported by a patient.

Conditions are recorded in different sources and levels of standardization. For example:

- Medical claims data include ICD-9-CM diagnosis codes that are submitted as part of a claim for health services and procedures.
- EHRs may capture a person's conditions in the form of diagnosis codes and symptoms as ICD-9-CM codes, but may not have a way to capture out-of-system conditions.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
condition_occurrence_id | Yes* | Integer | A unique identifier for each condition occurrence event. | This is not a value found in the EHR. Sites may choose to use a sequential value for this field
person_id | Yes* | Integer | A foreign key identifier to the person who is experiencing the condition. The demographic details of that person are stored in the person table.
condition_concept_id | Yes | Integer | A foreign key that refers to a standard condition concept identifier in the Vocabulary. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id ='SNOMED')</p> <p>select \* from concept where vocabulary_id ='SNOMED'  yields ~440,000 valid concept_ids.</p> If none are correct, use concept_id = 0
condition_start_date | Yes | Date| The date when the instance of the condition is recorded. | No date shifting.  
condition_end_date | No | Date| The date when the instance of the condition is considered to have ended | <p>No date shifting.</p> If this information is not available, set to NULL. 
condition_start_time | Yes | Datetime | The date and time when the instance of the condition is recorded. | No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
condition_end_time | No | Datetime | The date and time when the instance of the condition is considered to have ended | <p>No date shifting.</p> If this information is not available, set to NULL.  Full date and time. If there is no time associated with the date assert midnight.
condition_type_concept_id | Yes | Integer | A foreign key to the predefined concept identifier in the Vocabulary reflecting the source data from which the condition was recorded, the level of standardization, and the type of occurrence. For example, conditions may be defined as primary or secondary diagnoses, problem lists and person statuses. | <p>Please include valid concept ids (consistent with OMOP CDMv4). Predefined value set (valid concept_ids found in CONCEPT table where domain_id ='Condition Type')</p> <p>select \* from concept where domain_id ='Condition Type' yields 98 valid concept_ids.</p> <p>If none are correct, use concept_id = 0</p> If data source only identifies conditions as primary or secondary with no sequence number, use the following concept_ids: <ul><li>Inpatient primary: concept_id = 38000199</li> <li>Inpatient secondary: concept_id = 38000201</li> <li>Outpatient primary: concept_id = 38000230</li> <li>Outpatient secondary: concept_id = 38000231</li></ul>
stop_reason | No | Varchar | The reason, if available, that the condition was no longer recorded, as indicated in the source data. | <p>Valid values include discharged, resolved, etc. Note that a stop_reason does not necessarily imply that the condition is no longer occurring, and therefore does not mandate that the end date be assigned.</p>
provider_id | No | Integer | A foreign key to the provider in the provider table who was responsible for determining (diagnosing) the condition. | **In PEDSnet CDM v1, this field was previously called associated_provider_id**<p>Any valid provider_id allowed (see definition of providers in PROVIDER table)</p> Make a best-guess and document method used. Or leave blank
visit_occurrence_id | No* | Integer | A foreign key to the visit in the visit table during which the condition was determined (diagnosed).
condition_source_value | No* | Varchar | The source code for the condition as it appears in the source data. This code is mapped to a standard condition concept in the Vocabulary and the original code is, stored here for reference. | Condition source codes are typically ICD-9-CM diagnosis codes from medical claims or discharge status/visit diagnosis codes from EHRs. Use source_to_concept maps to translation from source codes to OMOP concept_ids.
condition_source_concept_id | No* | Integer | A foreign key to a condition concept that refers to the code used in the source| For example, if the condition is "Acute myeloid leukemia, without mention of having achieved remission" which has an icd9 code of 205.00 the condition source concept id is 44826430 which is the icd9 code concept that corresponds to the diagnosis 205.00.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.7.1 Additional Notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is \*\*NOT\*\* applied to condition_occurrence. All conditions are included for an active patient. For PEDSnet CDM V2, we limit condition_occurrences to final diagnoses only (not reason-for-visit and provisional surgical diagnoses such as those recored in EPIC OPTIME). In EPIC, final diagnoses includes both encounter diagnoses and billing diagnoses, problem lists (all problems, not filtered on "chronic" versus "provisional" unless local practices use this flag as intended). Medical History diagnosis are optional.
- Condition records are inferred from diagnostic codes recorded in the source data by a clinician or abstractionist for a specific visit. In the current version of the CDM, diagnoses extracted from unstructured data (such as notes) are not included.
- Source code systems, like ICD-9-CM, ICD-10-CM, etc., provide coverage of conditions. However, if the code does not define a condition, but rather is an observation or a procedure, then such information is not stored in the CONDITION_OCCURRENCE table, but in the respective tables instead. An example are ICD-9-CM procedure codes. For example, OMOP source-to-concept table uses the MAPPING_TYPE column to distinguish ICD9 codes that represent procedures rather than conditions.
- Condition source values are mapped to standard concepts for cflowonditions in the Vocabulary. Since the icd9-cm diagnosis codes are notB in the concept table, use the source_to_concept_map table where the icd9_code = source_code and the source_vocabulary_id =ICD9CM (icd_9) and target_vocabulatory_id=SNOMED (snomed-ct) to locate the correct condition_concept_id value.
- When the source code cannot be translated into a Standard Concept, a CONDITION_OCCURRENCE entry is stored with only the corresponding source_value and a condition_concept_id of 0.
- Codes written in the process of establishing the diagnosis, such as "question of" of and "rule out", are not represented here.

## 1.8 PROCEDURE_OCCURRENCE

The procedure occurrence domain contains records of significant activities or processes ordered by and/or carried out by a healthcare provider on the patient to have a diagnostic and/or therapeutic purpose that are not fully captured in another table (e.g. drug_exposure).

Procedures records are extracted from structured data in Electronic Health Records that capture source procedure codes using CPT-4, ICD-9-CM (Procedures), HCPCS or OPCS-4 procedures as orders.

**Only instantiated procedures are included in this table. Please exclude cancelled procedures**

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
procedure_occurrence_id | Yes* | Integer | A system-generated unique identifier for each procedure occurrence | This is not a value found in the EHR. Sites may choose to use a sequential value for this field
person_id | Yes* | Integer | A foreign key identifier to the person who is subjected to the procedure. The demographic details of that person are stored in the person table.
procedure_concept_id | Yes | Integer | A foreign key that refers to a standard procedure concept identifier in the Vocabulary. | <p>Valid Procedure Concepts belong to the "Procedure" domain. Procedure Concepts are based on a variety of vocabularies: SNOMED-CT (vocabulary_id ='SNOMED'), ICD-9-Procedures (vocabulary_id ='ICD9Proc'), CPT-4 (vocabulary_id ='CPT4' ), and HCPCS (vocabulary_id ='HCPCS')</p> <p>Procedures are expected to be carried out within one day. If they stretch over a number of days, such as artificial respiration, usually only the initiation is reported as a procedure (CPT-4 "Intubation, endotracheal, emergency procedure").</p> Procedures could involve the administration of a drug, in which case the procedure is recorded in the procedure table and simultaneously the administered drug in the drug table.
modifier_concept_id | No | Integer | A foreign key to a standard concept identifier for a modifier to the procedure (e.g. bilateral) |  <p>Valid Modifier Concepts belong to the "Modifier" domain. </p>
quantity | No |Integer |The quantity of procedures ordered or administered.
procedure_date | Yes | Date | The date on which the procedure was performed. 
procedure_time | Yes | Date | The date and time on which the procedure was performed. If there is no time associated with the date assert midnight. | **This field is a custom PEDSnet field**
procedure_type_concept_id | Yes | Integer | A foreign key to the predefined concept identifier in the Vocabulary reflecting the type of source data from which the procedure record is derived. (OMOP vocabulary_id = 'Procedure Type') | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id = 'Procedure Type')</p> <p>select \* from concept where vocabulary_id ='Procedure Type' yields 93 valid concept_ids.</p> If none are correct, use concept_id = 0
provider_id | No | Integer | A foreign key to the provider in the provider table who was responsible for carrying out the procedure. | <p>Any valid provider_id allowed (see definition of providers in PROVIDER table)</p> Document how selection was made.
visit_occurrence_id | No | Integer | A foreign key to the visit in the visit table during which the procedure was carried out. | See VISIT.visit_occurrence_id (primary key)
procedure_source_value | No* | Varchar | The source code for the procedure as it appears in the source data. This code is mapped to a standard procedure concept in the Vocabulary and the original code is stored here for reference. | Procedure_source_value codes are typically ICD-9, ICD-10 Proc, CPT-4, HCPCS, or OPCS-4 codes. All of these codes are acceptable source values.
procedure_source_concept_id | No* | Integer | A foreign key to a procedure concept that refers to the code used in the source. For example, if the procedure is "Anesthesia for procedures on eye; lens surgery" in the source which has a concept code in the vocabulary that is 2100658. The procedure source concept id will be 2100658.
modifier_source_value | No | Varchar | The source code for the modifier as it appears in the source data.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.8.1 Additional notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is \*\*NOT\*\* applied to procedure_occurrence. All procedures are included for an active patient. For PEDSnet CDM V2, we limit procedures_occurrences to billing procedures only (not surgical diagnoses).
- Procedure Concepts are based on a variety of vocabularies: SNOMED-CT, ICD-9-Proc, CPT-4, HCPCS and OPCS-4.
- Procedures could reflect the administration of a drug, in which case the procedure is recorded in the procedure table and simultaneously the administered drug in the drug table.
- The Visit during which the procedure was performed is recorded through a reference to the VISIT_OCCURRENCE table. This information is not always available.
- The Provider carrying out the procedure is recorded through a reference to the PROVIDER table. This information is not always available.

## 1.9 OBSERVATION

The observation domain captures clinical facts about a patient obtained in the context of examination, questioning or a procedure. The observation domain supports capture of data not represented by other domains, including unstructored measurements, medical history and family history. For the PEDSnet CDM version 2, the observations listed below are extracted from source data. Please assign the specific concept_ids listed in the table below to these observations as observation_concept_ids. Non-standard PCORnet concepts require concepts that have been entered into an OMOP-generated vocabulary (OMOP provided vocabulary_id ='PCORNet').

NOTE: DRG and DRG Type require special logic/processing described below.

- Biobank availability
- Admitting source
- Discharge disposition
- Discharge status
- Chart availability
- DRG (requires special logic - see Note 1 below)
- Tobacco Status (see Note 4)

Use the following table to populate observation_concept_ids for the observations listed above. The vocabulary column is used to highlight non-standard codes from vocabulary_id='Observation Type' and vocabulary_id='PCORNet' and one newly added standard concept from vocabulary_id='SNOMED'.

**Table 1: Observation concept IDs for PCORnet concepts. The vocabulary id 'PCORNet' contains concept specific to PCORNet requirements and standards.**

Concept Name | Observation concept ID | Vocab ID | Value as concept ID | Concept description | Vocab ID
 --- | --- | --- | --- | --- | ---
Biobank flag (see Note 2) | 4001345 | | 4188539 | Yes
Biobank flag | 4001345 | | 4188540 | No
Biobank flag | 4001345 | | 44814650 | No information |PCORNet
Biobank flag | 4001345 | | 44814653 | Unknown | PCORNet
Biobank flag | 4001345 | | 44814649 | Other | PCORNet
Admitting source | 4145666 | | 44814670 | Adult Foster Home
Admitting source | 4145666 | | 44814671 | Assisted Living Facility
Admitting source | 4145666 | | 44814672 | Ambulatory Visit
Admitting source | 4145666 | | 8870 | Emergency Department
Admitting source | 4145666 | | 44814674 | Home Health
Admitting source | 4145666 | | 44814675 | Home / Self Care
Admitting source | 4145666 | | 8546 | Hospice
Admitting source | 4145666 | | 38004279 | Other Acute Inpatient Hospital
Admitting source | 4145666 | | 44814678 | Nursing Home (Includes ICF)
Admitting source | 4145666 | | 44814679 | Rehabilitation Facility
Admitting source | 4145666 | | 44814680 | Residential Facility | PCORNet
Admitting source | 4145666 | | 8863 | Skilled Nursing Facility
Admitting source | 4145666 | | 44814650 | No information | PCORNet
Admitting source | 4145666 | | 44814653 | Unknown | PCORNet
Admitting source | 4145666 | | 44814649 | Other | PCORNet
Discharge disposition (See Note 3) | 44813951 | SNOMED | 4161979 | Discharged alive
Discharge disposition | 44813951 | SNOMED | 4216643 | Expired
Discharge disposition | 44813951 | SNOMED | 44814650 | No information | PCORNet
Discharge disposition | 44813951 | SNOMED | 44814653 | Unknown | PCORNet
Discharge disposition | 44813951 | SNOMED | 44814649 | Other | PCORNet
Discharge status (see Note 3) | 4137274 | | 38004205 | Adult Foster Home
Discharge status | 4137274 | | 38004301 | Assisted Living Facility
Discharge status | 4137274 | | 4021968 | Against Medical Advice
Discharge status | 4137274 | | 44814693 | Absent without leave | PCORNet
Discharge status | 4137274 | | 4216643 | Expired
Discharge status | 4137274 | | 38004195 | Home Health
Discharge status | 4137274 | | 8536 | Home / Self Care
Discharge status | 4137274 | | 8546 | Hospice
Discharge status | 4137274 | | 38004279 | Other Acute Inpatient Hospital
Discharge status | 4137274 | | 8676 | Nursing Home (Includes ICF)
Discharge status | 4137274 | | 8920 | Rehabilitation Facility
Discharge status | 4137274 | | 44814701 | Residential Facility | PCORNet
Discharge status | 4137274 | | 8717 | Still In Hospital
Discharge status | 4137274 | | 8863 | Skilled Nursing Facility
Discharge status | 4137274 | | 44814653 | Unknown | PCORNet
Discharge status | 4137274 | | 44814649 | Other | PCORNet
Discharge status | 4137274 | | 44814650 | No information | PCORNet
Chart availability (See Note 2) | 4030450 | | 4188539 | Yes
Chart availability | 4030450 | | 4188540 | No
Chart availability  | 4030450 | | 44814650 | No information |PCORNet
Chart availability  | 4030450 | | 44814653 | Unknown | PCORNet
Chart availability  | 4030450 | | 44814649 | Other | PCORNet
Tobacco Use |*concept id pending* | |*concept id pending* | Current every day smoker
Tobacco Use | *concept id pending*| |*concept id pending* |  Current some day smoker
Tobacco Use |*concept id pending* | |*concept id pending* |  Former smoker
Tobacco Use |*concept id pending* | |*concept id pending* | Never Smoker
Tobacco Use |*concept id pending* | |*concept id pending* |Smoker, current status unknown
Tobacco Use |*concept id pending* | |*concept id pending* | Unknown if ever smoked
Tobacco Use |*concept id pending* | |*concept id pending* | Heavy tobacco smoker
Tobacco Use |*concept id pending* | |*concept id pending* | Light tobacco smoker
Tobacco Use |*concept id pending* | |44814650 |No information | PCORNet
Tobacco Use |*concept id pending* | |44814653| Unknown| PCORNet
Tobacco Use |*concept id pending* | |44814649| Other| PCORNet
Tobacco Type |*concept id pending* | |*concept id pending* | Cigarettes only
Tobacco Type |*concept id pending* | |*concept id pending* | Other tobacco only
Tobacco Type |*concept id pending* | |*concept id pending* | Cigarettes and other tobacco
Tobacco Type |*concept id pending* | |*concept id pending* | None
Tobacco Type |*concept id pending* | | 44814650 |No information | PCORNet
Tobacco Type |*concept id pending* | |44814653| Unknown| PCORNet
Tobacco Type |*concept id pending* | |44814649| Other| PCORNet

**Note 1**: For DRG, use the following logic (must use vocabulary version 5):

- The DRG value must be three digits as text. Put into value_as_string in observation
- For all DRGs, set observation_concept_id = 3040464 (hospital discharge DRG)
- To obtain correct value_as_concept_id for the DRG:
    - If the date for the DRG \< 10/1/2007, use concept_class_id = "DRG", invalid_date = "9/30/2007", invalid_reason = 'D' and the DRG value=CONCEPT.concept_code to query the CONCEPT table for correct concept_id to use as value_as_concept_id.
    - If the date for the DRG \>=10/1/2007, use concept_class_id = "MS-DRG", invalid_reason = NULL and the DRG value = CONCEPT.concept_code to query the CONCEPT table for the correct concept_id to use as value_as_concept_id.

**Note 2**: In the Observation table, the biobank flag and chart availability concept_ids can appear multiple times capturing changes in patient consent over time. The temporally most recent observation will be used to determine the current consent status.

**Note 3:** Discharge disposition and discharge status appear only once per visit_occurence. These vales can change across different visit_occurrences. Use the visit_occurrence_id to tie these observations to the corresponding visit.

**Note 4:** If tobacco information is available at the visit level, please provide this information. If it is not sites are welcomed to make a high level assertion about tobacco use and tobacco type information for individuals in the cohort.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
observation_id | Yes* | Integer |  A unique identifier for each observation. | This is not a value found in the EHR. Sites may choose to use a sequential value for this field
person_id | Yes* | Integer | A foreign key identifier to the person about whom the observation was recorded. The demographic details of that person are stored in the person table.|
observation_concept_id | Yes | Integer | A foreign key to the standard observation concept identifier in the Vocabulary. | Lab results and vitals are not stored in this table in V5 but are stored in the Measurement table.
observation_date | Yes* | Date | The date of the observation. | No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
observation_time | No | Datetime | The time of the observation. | No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
observation_type_concept_id | Yes | Integer | A foreign key to the predefined concept identifier in the Vocabulary reflecting the type of the observation. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id ='Observation Type')</p> <p>select \* from concept where vocabulary_id = 'Observation Type' yields 11 valid concept_ids.</p> FOR PEDSnet CDM V2, all of our observations are coming from electronic health records so *set this field to concept_id = 38000280* (observation recorded from EMR). When we get data from patients, we will include the concept_id = 44814721
value_as_number | No\* (see convention) | Float | The observation result stored as a number. This is applicable to observations where the result is expressed as a numeric value. | Value must be represented as at least one of {value_as_number, value_as_string or values_as_concept_id}. There are a few exceptions in vocabulary id 99 where all three value_as_\* fields are NULL.
value_as_string | No\* (see convention) | Varchar | The observation result stored as a string. This is applicable to observations where the result is expressed as verbatim text. | Value must be represented as at least one of {value_as_number, value_as_string or values_as_concept_id}. There are a few exceptions in vocabulary id PCORNet where all three value_as_\* fields are NULL.
value_as_concept_id | No\* (see convention) | Integer | A foreign key to an observation result stored as a concept identifier. This is applicable to observations where the result can be expressed as a standard concept from the Vocabulary (e.g., positive/negative, present/absent, low/high, etc.). | Value must be represented as at least one of {value_as_number, value_as_string or values_as_concept_id}. There are a few exceptions where all three value_as_\* fields are NULL.
qualifier_concept_id | No | Integer | A foreign key to standard concept identifier for a qualifier (e.g severity of drug-drug interaction alert) | <p>Predefined value set (valid concept_ids found in CONCEPT table where domain_id='Observation' and concept_class_id ='Qualifier Value')</p> <p>select \* from concept where domain_id='Observation' and concept_class_id ='Qualifier Value' yields 10912 valid concept_ids.</p> <p>If none are correct, use concept_id = 0.</p>
unit_concept_id | No | Integer | A foreign key to a standard concept identifier of measurement units in the Vocabulary. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where domain_id='Unit' and vocabulary_id ='UCUM')</p> <p>select \* from concept where domain_id='Unit' and vocabulary_id ='UCUM' yields 920 valid concept_ids.</p> <p>If none are correct, use concept_id = 0.</p>
provider_id | No | Integer | A foreign key to the provider in the provider table who was responsible for making the observation.
visit_occurrence_id | No* | Integer | A foreign key to the visit in the visit table during which the observation was recorded.
observation_source_value | No | Varchar | The observation code as it appears in the source data. This code is mapped to a standard concept in the Vocabulary and the original code is, stored here for reference.
observation_source_concept_id| No |Integer | A foreign key to a concept that refers to the code used in the source. |
unit_source_value | No | Integer | The source code for the unit as it appears in the source data. This code is mapped to a standard unit concept in the Vocabulary and the original code is, stored here for reference.
qualifier_source_value |No | Varchar | The source value associated with a qualifier to characterize the observation

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.9.1 Additional Notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is \*\*NOT\*\* applied to observations. All observations are included for an active patient. For PEDSnet CDM V2, we limit observations to only those that appear in Table 1.
- Observations have a value represented by one of a concept ID, a string, \*\*OR\*\* a numeric value.
- The Visit during which the observation was made is recorded through a reference to the VISIT_OCCURRENCE table. This information is not always available.
- The Provider making the observation is recorded through a reference to the PROVIDER table. This information is not always available.
- Observations obtained using standardized methods (e.g. laboratory assays) that produce discrete results are recorded by preference in the MEASUREMENT table.

## 1.10 OBSERVATION PERIOD

The observation period domain is designed to capture the time intervals in which data are being recorded for the person. An observation period is the span of time when a person is expected to have a clinical fact represented in the PEDSNet version 2 data modelo. This table is used to generate the PCORnet CDM enrollment table.

While analytic methods can be used to calculate gaps in observation periods that will generate multiple records (observation periods) per person, for PEDSnet, the logic has been simplified to generate a single observation period row for each patient.

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
Observation_period_id | Yes | Integer | A system-generate unique identifier for each observation period | This is not a value found in the EHR. Sites may choose to use a sequential value for this field.
person_id | Yes* | Integer | A foreign key identifier to the person who is experiencing the condition. The demographic details of that person are stored in the person table.
Observation_period_start_date | Yes | Date | The start date of the observation period for which data are available from the data source | <p>Use the earliest clinical fact date available for this patient.</p> No date shifting.  
Observation_period_end_date | No* | Date | The end date of the observation period for which data are available from the source. | <p>Use the latest clinical fact date available for this patient. If there exists one or more records in the DEATH table for this patient, use the latest date recorded in that table.</p> For patients who are still in the hospital or ED or other facility at the time of data extraction, leave this field NULL. 
Observation_period_start_time | Yes* | Datetime | The start date of the observation period for which data are available from the data source | <p>Use the earliest clinical fact time available for this patient.</p> No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
Observation_period_end_time | No | Datetime | The end date of the observation period for which data are available from the source. | <p>Use the latest clinical fact time available for this patient. If there exists one or more records in the DEATH table for this patient, use the latest date recorded in that table.</p> For patients who are still in the hospital or ED or other facility at the time of data extraction, leave this field NULL.  Full date and time. If there is no time associated with the date assert midnight.

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.10.1 Additional Notes

- Because the 1/1/2009 date limitation for "active patients" is not used to limit visit_occurrance, the start_date of an observation period for an active PEDSnet patient may be prior to 1/1/ 2009.


## 1.11 DRUG EXPOSURE (\*\*DRAFT\*\*)

The drug exposure domain captures any biochemical substance that is introduced in any way to a patient. This can be evidence of prescribed, over the counter, administered (IV, intramuscular, etc), immunizations or dispensed medications. These events could be linked to procedures or encounters where they are administered or associated as a result of the encounter.

EHRs may store medications in different vocabularies (GPI,NDC etc). 

Exclusions:

1. Cancelled Medication Orders
2. Missed Medication administrations

**Note 1**: The effective_dose_drug is the dose basis.(Eg. 45 mg/kg/dose). This is the discrete dose value from the source data if available. If the discrete dose value is **not** available from the source data, then compute the dose basis by looking for a weight observation **+/- 60 days of the date of the medication**. (Eg. Total Amount/**(divided by)**Weight) (Dose per kg)

**Note 2**: The quantity is the actual dose given. (Eg. 450 mg for 10 kg patient)

**Note 3**: The dose_unit_concept_id is the unit of the effective dose.

**Note 4**: For dispensing records, compute the dose basis by looking for a weight observation +/- 60 days of the dispensed date.

**Note 5:** For the sig, encode the value using XML. 

<ul> <li> Element 1: Actual SIG from source data </li> <li> Element 2: Raw "Supply/Quantity" (Examples: "1 bottle" "10 ml Bottle" "1 pack"</li> <li>Element 3: Refills</li></ul>

```xml
    <XML>
    <SIG>1/2 capful in 4 oz clear liquid</SIG>
    <QUANTITY>1 jar</QUANTITY>
    <REFILLS>2</REFILLS>
    </XML>
```


Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
 drug_exposure_id | Yes | Integer | A system-generated unique identifier for each drug exposure | This is not a value found in the EHR. Sites may choose to use a sequential value for this field.
person_id | Yes* |Integer | A foreign key identifier to the person who is experiencing the condition. The demographic details of that person are stored in the person table.
drug_concept_id| Yes* | Integer | A foreign key that refers to a standard drug concept identifier in the Vocabulary. | Valid drug concept IDs are mapped to RxNorm using the source to concept map table to transform source codes (GPI, NDC etc to the RxNorm target)
drug_exposure_start_date| Yes* | Date |The start date of the utilization of the drug. The start date of the prescription, the date the prescription was filled, the date a drug was dispensed or the date on which a drug administration procedure was recorded are acceptable. | No date shifting. 
drug_exposure_end_date| No* |Date | The end date of the utilization of the drug | No date shifting.
drug_exposure_start_time| Yes | Datetime |The start date and time of the utilization of the drug. The start date of the prescription, the date the prescription was filled, the date a drug was dispensed or the date on which a drug administration procedure was recorded are acceptable. | No date shifting. Full date and time. If there is no time associated with the date assert midnight.
drug_exposure_end_time| No |Datetime | The end date and time of the utilization of the drug | No date shifting. Full date and time. If there is no time associated with the date assert midnight.
    drug_type_concept_id| Yes | Integer | A foreign key to a standard concept identifier of the type of drug exposure in the Vocabulary as represeneted in the source data | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where domain_id ='Drug Type')</p> <p>select \* from concept where domain_id ='Drug Type' yields 12 valid concept_ids.</p> <p>If none are correct, use concept_id = 0.</p> For the PEDSnet observation listed above, use the following concept_ids: <ul><li>Prescription dispensed in pharmacy (dispensed meds pharma information): concept_id = 38000175</li> <li>Inpatient administration (MAR entries): concept_id = 38000180</li> <li>Prescription written: concept_id = 38000177</li></ul>
stop_reason| No | Varchar | The reason, if available, where the medication was stopped, as indicated in the source data. | <p>Valid values include therapy completed, changed, removed, side effects, etc. Note that a stop_reason does not necessarily imply that the medication is no longer being used at all, and therefore does not mandate that the end date be assigned.</p>
refills| No | Integer | The number of refills after the initial prescrition||
quantity| No | Integer | The quantity of the drugs as recorded in the original prescription or dispensing record| See Note 2|
days_supply| No | Integer | The number of days of supply the meidcation as recorded in the original prescription or dispensing record||
sig| No | CLOB (XML Structure) | The directions on the drug prescription as recorded in the original prescription (and printed on the container) or the dispensing record| See Note 5|
route_concept_id| No | Integer | A foreign key that refers to a standard administration route concept identifier in the Vocabulary. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where concept_class_id='Dose Form')</p> <p>select * from omop5.concept where concept_class_id='Dose Form' yields 357 valid concept_ids.</p> <p>If none are correct, use concept_id = 0.</p>
effective_drug_dose| No | Float | Numerical value of drug dose for this drug_exposure record| See note 1|
dose_unit_concept_id| No | Integer | A foreign key to a predefined concept in the Standard Vocabularies reflecting the unit the effective drug_dose value is expressed|See note 3|
lot_number| No | Varchar | An identifier to determine where the product originated||
provider_id| No | Integer | A foreign key to the provider in the provider table who initiated (prescribed) the drug exposure |<p>Any valid provider_id allowed (see definition of providers in PROVIDER table)</p> Document how selection was made.
visit_occurrence_id| No | Integer | A foreign key to the visit in the visit table during which the drug exposure initiated. | See VISIT.visit_occurrence_id (primary key)
drug_source_value| No*| Varchar | The source drug value as it appears in the source data. The source is mapped to a standard RxNorm concept and the original code is stored here for reference.
drug_source_concept_id| No | Integer | A foreign key to a drug concept that refers to the code used in the source | In this case, if you are transforming drugs from GPI or NDC to RXNorm. The concept id that corresponds to the GPI or NDC value for the drug belongs here.
route_source_value| No| Varchar |The information about the route of administration as detailed in the source ||
dose_unit_source_value| No| Varchar | The information about the dose unit as detailed in the source ||

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.11.1 Additional Notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is \*\*NOT\*\* applied to drug exposures. All drug exposures are included for an active patient. 
- The Visit during which the drug exposure was initiated by is recorded through a reference to the VISIT_OCCURRENCE table. This information is not always available.
- The Provider initating th drug exposure is recorded through a reference to the PROVIDER table. This information is not always available.

## 1.12 MEASUREMENT (\*\*DRAFT\*\*)

The measurement domain captures measurement orders and measurement results. The measurement domain can contain laboratory results and vital signs.

Specifically this table includes:
- Height/length in cm (use numeric precision as recorded in EHR)
- Height/length type
- Weight in kg (use numeric precision as recorded in EHR)
- Body mass index in kg/m<sup>2</sup> (extracted only if height and weight are not present)
- Systolic blood pressure in mmHg
    - Where multiple readings are present on the same encounter, create observation records for \*\*ALL\*\* readings
- Diastolic blood pressure in mmHg
    - Where multiple readings are present on the same encounter, create observation records for \*\*ALL\*\* readings
- Blood pressure position is described by the selection of a concept_id that contains the BP position as describe below. For example, in Table 1, concept_id 3018586 is Systolic Blood Pressure, Sitting. This concept_id identifies both the measurement (Systolic BP) and the BP position (sitting).
- Vital source
- Component Level Labs **(Network list pending)**

**Table 3: Measurement concept IDs for PCORnet concepts. Concept_ids from vocabulary_id 99 are non-standard codes.**

Concept Name | Observation concept ID | Vocab ID | Value as concept ID | Concept description | Vocab ID
 --- | --- | --- | --- | --- | ---
Vital | 3013762 | | See Note 1 | Weight
Vital | 3023540 | | See Note 1 | Height
Vital | 3034703 | | See Note 2 | Diastolic Blood Pressure - Sitting
Vital | 3019962 | | See Note 2 | Diastolic Blood Pressure - Standing
Vital | 3013940 | | See Note 2 | Diastolic Blood Pressure - Supine
Vital | 3012888 | | See Note 2 | Diastolic BP Unknown/Other
Vital | 3018586 | | See Note 2 | Systolic Blood Pressure - Sitting
Vital | 3035856 | | See Note 2 | Systolic Blood Pressure - Standing
Vital | 3009395 | | See Note 2 | Systolic Blood Pressure - Supine
Vital | 3004249 | | See Note 2 | Systolic BP Unknown/Other
Vital | 3038553 | | See Note 1 | BMI
Measurement Type | 44818704 | Measurement Type | See Note 3 | Patient reported
Measurement Type | 44818701 | Measurement Type | See Note 3 | Vital sign

**Note 1**: For height, weight and BMI observations, insert the recorded measurement into the value_as_number field.
**Note 2**: Systolic and diastolic pressure measurements will generate two observation records one for storing the systolic blood pressure measurement and a second for storing the diastolic blood pressure measurement. Select the right SBP or DBP concept code that also represents the CORRECT recording position (supine, sitting, standing, other/unknown). To tie the two measurements together (the systolic BP measurement and the diastolic BP measurement records), use the FACT_RELATIONSHIP table.

Example: Person_id = 12345 on visit_occurrence_id = 678910 had orthostatic blood pressure measurements performed as follows:

- Supine: Systolic BP 120; Diastolic BP 60
- Standing: Systolic BP 144; Diastolic BP 72

Four rows will be inserted into the observation table. Showing only the relevant columns:

Observation_id | Person_id | Visit_occurrence_id | measurement_concept_id | measurement_type_concept_id | Value_as_Number | Value_as_String | Value_as_Concept_ID
 --- | --- | --- | --- | --- | --- | --- | ---
66661 | 12345 | 678910 | 3009395 | 38000280 | 120 | |
66662 | 12345 | 678910 | 3013940 | 38000280 | 60 | |
66663 | 12345 | 678910 | 3035856 | 38000280 | 144 | |
66664 | 12345 | 678910 | 3019962 | 38000280 | 72 | |

- Measurement_concept_id = 3009395 = systolic BP - supine; measurement_concept_id = 3013940 = diastolic BP supine
- Measurement_concept_id = 3035856 = systolic BP standing; measurement_concept_id = 3019962 = diastolic BP standing
- measurement_type_concept_id = 38000280 (observation recorded from EMR).

To link these two values, use the fact relationship table:

Domain_concept_id_1 | fact_id_1 | Domain_concept_id_2 | fact_id_2 | relationship_concept_id
--- | --- | --- | --- | ---
Measurement | 66661 | Measurement | 66662 |  Asso with finding 
Measurement | 66662 | Measurement | 66661 |  Asso with finding 
Measurement | 66663 | Measurement | 66664 |  Asso with finding 
Measurement | 666624 | Measurement | 66663 |  Asso with finding 

Because the domain concept id and relationship concept id are integers the following is an example of how this data will be represented:

Domain_concept_id_1 | fact_id_1 | Domain_concept_id_2 | fact_id_2 | relationship_concept_id
--- | --- | --- | --- | ---
21 | 66661 | 21 | 66662 |  44818792 
21 | 66662 | 21 | 66661 |  44818792 
21 | 66663 | 21 | 66664 |  44818792 
21 | 666624 | 21 | 66663 | 44818792


- Two rows in the FACT_RELATIONSHIP table link the *supine*  diastolic BP to the supine systolic BP.
- Two rows in the FACT_RELATIONSHIP table link the *standing* diastolic BP to the standing systolic BP.

**Note 3**: Measurement type concept_ids are used as values for the measurement_type_concept_id field.
In addition, the following observations are derived via the DCC (concept_ids to be assigned in future version of this document. However, concept_ids are not needed for ETL since these observations will be derived/calculated using scripts developed by DCC):

- Body mass index in kg/m<sup>2</sup> if not directly extracted
- Height/length z score for age/sex using NHANES 2000 norms for measurements at which the person was <240 months of age. In the absence of a height/length type for the measurement, recumbent length is assumed for ages \<24 months, and standing height thereafter.
- Weight z score for age/sex using NHANES 2000 norms for measurements at which the person was \<240 months of age.
- BMI z score for age/sex using NHANES 2000 norms for visits at which the person was between 20 and 240 months of age.
- Systolic BP z score for age/sex/height using NHBPEP task force fourth report norms.
- Diastolic BP z score for age/sex/height using NHBPEP task force fourth report norms.

Exclusions:

1. Cancelled Lab orders
2. Lab orders that are 'NOT DONE' or 'INCOMPLETE'

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
measurement_id | Yes* | Integer | A system-generated unique identifier for each measurement | This is not a value found in the EHR. Sites may choose to use a sequential value for this field.
person_id | Yes* | Integer | A foreign key identifier to the person who the measurement is being documented for. The demographic details of that person are stored in the person table.
measurement_concept_id | Yes* | Integer | A foreign key to the standard measurement concept identifier in the Vocabulary. | <p>Valid Measurement Concepts belong to the "Measurement" domain. Measurement Concepts are based mostly on the LOINC vocabulary, with some additions from SNOMED-CT.</p> <p>Measurement must have an object represented as a concept, and a finding. A finding (see below) is represented as a concept, a numerical value or a verbatim string or more than one of these.</p> <p>There are three Standard Vocabularies defined for measurements:</p> <p>Laboratory tests and values: Logical Observation Identifiers Names and Codes (**LOINC**) (Vocabulary_id=LOINC).</p> <p>(FYI: Regenstrief also maintains the **"LOINC Multidimensional Classification"** Vocabulary_id=LOINC Hierarchy)</p> <p>Qualitative lab results: A set of SNOMED-CT Qualifier Value concepts (vocabulary_id=SNOMED)</p> <p>Laboratory units: Unified Code for Units of Measure (**UCUM**( )Vocabulary_id=UCUM)</p> <p>All other findings and observables: SNOMED-CT (Vocabulary_id=SNOMED).</p> For vital signs, pull information from flow sheet rows (EPIC sites only)
measurement_date| Yes* | Date | The date of the measurement. | No date shifting.  
measurement_time| No| Datetime | The time of the measurement. | No date shifting.  Full date and time. If there is no time associated with the date assert midnight.
measurement_type_concept_id | Yes | Integer | A foreign key to the predefined concept identifier in the Vocabulary reflecting the type of the measurement. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id =Meas Type)</p> <p>select \* from concept where vocabulary_id =Meas Type yields 5 valid concept_ids.</p> For Pedsnet CDM v2, please use the following: <ul><li>Vital Sign= 44818701</li><li>Lab result =  44818702</li><li>Pathology finding = 44818703</li><li>Patient reported value = 44818704</li> <li> Derived Value = 45754907</li></ul> 
operator_concept_id| No| Integer | A foreign key identifier to the mathematical operator that is applied to the value_as_number.Operators are <, ≤, =, ≥, >| Valid operator concept id are found in the concept table <p> select \* from concept where domain_id='Meas Value Operator' yields 5 valid concept ids. <ul> <li> Operator <= : 4171754 </li> <li> Operator >= : 4171755      </li> <li> Operator < : 4171756 </li> <li> Operator =   4172703 </li> <li> Operator > : 4172704 </li> </ul>|
value_as_number | No\* (see convention) | Float | The measurement result stored as a number. This is applicable to measurements where the result is expressed as a numeric value. | Value must be represented as at least one of {value_as_number, value_as_string or values_as_concept_id}.
value_as_concept_id | No\* (see convention) | Integer | A foreign key to an observation result stored as a concept identifier. This is applicable to observations where the result can be expressed as a standard concept from the Vocabulary (e.g., positive/negative, present/absent, low/high, etc.). | Value must be represented as at least one of {value_as_number, value_as_string or values_as_concept_id}. Valid concepts are found in the concept table <p> select \* from concept where domain_id='Meas Value' yiels 86 valid concept ids.</p>
unit_concept_id | No | Integer | A foreign key to a standard concept identifier of measurement units in the Vocabulary. | <p>Please include valid concept ids (consistent with OMOP CDMv5). Predefined value set (valid concept_ids found in CONCEPT table where vocabulary_id = UCUM)</p> <p>select \* from concept where vocabulary_id = UCUM yields 912 valid concept_ids.</p> <p>If none are correct, use concept_id = 0.</p> For the PEDSnet observation listed above, use the following concept_ids: <ul><li>Centimeters (cm): concept_id = 8582</li> <li>Kilograms (kg): concept_id = 9529</li> <li>Kilograms per square meter (kg/m<sup>2</sup>): concept_id = 9531</li> <li>Millimeters mercury (mmHG): concept_id = 8876</li></ul>
range_low | No | Float | <p>Optional - Do not transmit to DCC</p> The lower limit of the normal range of the measurement. It is not applicable if the observation results are non-numeric or categorical, and must be in the same units of measure as the measurement value.
range_high | No | Float | <p>Optional - Do not transmit to DCC.</p> The upper limit of the normal range of the measurement. It is not applicable if the observation results are non-numeric or categorical, and must be in the same units of measure as the measurement value.
provider_id | No | Integer | A foreign key to the provider in the provider table who was responsible for making the measurement.
visit_occurrence_id | No* | Integer | A foreign key to the visit in the visit table during which the observation was recorded.
measurement_source_value | Yes | Varchar | The measurement name as it appears in the source data. This code is mapped to a standard concept in the Standardized Vocabularies and the original code is, stored here for reference.| This is the name of the value as it appears in the source system. For lab values, it is suggested to include the (LAB ID\| PROCEDURE NAME \| COMPONENT NAME) combination as the measurement source value where applicable.
measurement_source_concept_id| No| Integer | A foreign key to a concept that refers to the code used in the source.| This is the concept id that maps to the source value in the standard vocabulary.
unit_source_value| No| Varchar | The source code for the unit as it appears in the source data. This code is mapped to a standard unit concept in the Standardized Vocabularies and the original code is, stored here for reference.| Raw unit value (Ounces,Inches etc)
value_source_value| Yes*| Varchar | The source value associated with the structured value stored as numeric or concept. This field can be used in instances where the source data are transformed|<ul> <li>For BP values include the raw 'systolic/diastolic' value Eg. 120/60</li><li>If there are transformed values (Eg. Weight and Height) please insert the raw data before transformation.</li></ul>

**\* This field is important for responding to PCORNet queries. If sites have any information on this filed in the source EHR then these fields should be populated in the PEDSnet CDM instance**

#### 1.12.1 Additional Notes

- The 1/1/2009 date limitation that is used to define a PEDSnet active patient is \*\*NOT\*\* applied to measurements. All measurements are included for an active patient. For PEDSnet CDM V1, we limit measurements to only those that appear in Table 3 (for vital signs).
- Measurements have a value represented by one of a concept ID, a string, \*\*OR\*\* a numeric value.
- The Visit during which the measurement was made is recorded through a reference to the VISIT_OCCURRENCE table. This information is not always available.
- The Provider making the measurement is recorded through a reference to the PROVIDER table. This information is not always available.

## 1.13 FACT RELATIONSHIP

The fact relationship domain contains details of the relationships between facts within one domain or across two domains, and the nature of the relationship. Examples of types of possible fact relationships include: person relationships (mother-child linkage), care site relationships (representing the hierarchical organization structure of facilities within health systems), drug exposures provided due to associated indicated condition, devices used during the course of an associated procedure, and measurements derived from an associated specimen. All relationships are directional, and each relationship is represented twice symmetrically within the fact relationship table. 

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
Domain_concept_id_1|Yes| Integer |	The concept representing the domain of fact one, from which the corresponding table can be inferred
Fact_id_1|	Yes | Integer |The unique identifier in the table corresponding to the domain of fact one.
Domain_concept_id_2|Yes| Integer |	The concept representing the domain of fact two, from which the corresponding table can be inferred.
Fact_id_2 |	Yes | Integer |	The unique identifier in the table corresponding to the domain of fact two.
Relationship_concept_id	|Yes| Integer |A foreign key to a standard concept identifier of relationship in the Standardized Vocabularies.

#### 1.13.1 Additional Notes
- Blood Pressure Systolic and Diastolic Blood Pressure Values will be mapped using the fact relationship table.

## 1.14 VISIT_PAYER

The visit payer table documents insurance information as it relates to a visit in visit_occurrence. For this reason the key of this table will be visit_occurrence_id and visit_payer_id. **This table is CUSTOM to Pedsnet.**

Field |Required | Data Type | Description | PEDSnet Conventions
 --- | --- | --- | --- | ---
visit_payer_id | Yes | Integer |A system-generated unique identifier for each visit payer relationship.
visit_occurrence_id | Yes | Integer | A foreign key to the visit in the visit table where the payer was billed for the visit.
plan_name | Yes | Varchar| The untransformed payer/plan name from the source data
plan_type | No | Varchar |  A standardized interpretation of the plan structure; proposed value set would be HMO, PPO, POS, Fee for service, Other/unknown
plan_class | Yes | Varchar | A list of the "payment sources" most often used in demographic analyses; proposed value would be Private/commercial, Medicaid/sCHIP, Medicare, Other public, Self-pay, Other/unknown

* * *

**APPENDIX**

**PEDSnet-specific is supported by OMOP-supported Vocabulary id=PCORNet, which contains all of the additional concept_id codes needed in PEDSnet for PCORnet CDM V1.0 and 2.0**

### A1. ABMS Specialty Category to OMOP V5 Specialty Mapping
http://www.abms.org/member-boards/specialty-subspecialty-certificates/

ABMS Specialty Category | OMOP Supported Concept for Provider ID | OMOP Concept_name | Domain_id | Vocabulary id
--- | --- | --- | --- | ---
Addiction Psychiatry |38004498 | Addiction Medicine | Provider Specialty | Specialty   
Adolescent Medicine |45756747 |Adolescent Medicine|Provider Specialty | ABMS
Adult Congenital Heart Disease |45756748  |Adult Congenital Heart Disease|Provider Specialty |  ABMS
Advanced Heart Failure and Transplant Cardiology |45756749 |Advanced Heart Failure and Transplant Cardiology|  Provider Specialty | ABMS
Aerospace Medicine |45756750 |Aerospace Medicine| Provider Specialty | ABMS
Allergy and Immunology | 38004448 | Allergy/Immunology                | Provider Specialty | Specialty    
Anesthesiology |38004450 | Anesthesiology   | Provider Specialty | Specialty    
Anesthesiology Critical Care Medicine |45756751 |Anesthesiology Critical Care Medicine|Provider Specialty | Specialty
Blood Banking/Transfusion Medicine |45756752|Blood Banking/Transfusion Medicine|Provider Specialty |  ABMS
Brain Injury Medicine |45756753 |Brain Injury Medicine|Provider Specialty |  ABMS
Cardiology |38004451 | Cardiology                       | Provider Specialty | Specialty 
Cardiovascular Disease |45756754 |Cardiovascular Disease|Provider Specialty |  ABMS
Child Abuse Pediatrics |45756755 |Child Abuse Pediatrics |Provider Specialty |  ABMS
Child and Adolescent Psychiatry |45756756|Child and Adolescent Psychiatry|Provider Specialty |  ABMS
Clinical Biochemical Genetics |45756757|Clinical Biochemical Genetics|Provider Specialty |  ABMS
Clinical Cardiac Electrophysiology |45756758|Clinical Cardiac Electrophysiology |Provider Specialty |  ABMS
Clinical Cytogenetics |45756759 |Clinical Cytogenetics|Provider Specialty |  ABMS
Clinical Genetics (MD) |45756760|Clinical Genetics (MD)|Provider Specialty |  ABMS
Clinical Informatics |45756761|Clinical Informatics|Provider Specialty |  ABMS
Clinical Molecular Genetics |45756762|Clinical Molecular Genetics|Provider Specialty |  ABMS
Clinical Neurophysiology |45756763|Clinical Neurophysiology |Provider Specialty |  ABMS
Colon and Rectal Surgery | 38004471 | Colorectal Surgery              | Provider Specialty | Specialty   
Complex General Surgical Oncology |45756764|Complex General Surgical Oncology| Provider Specialty | ABMS
Congenital Cardiac Surgery |45756765|Congenital Cardiac Surgery |Provider Specialty |  ABMS
Critical Care Medicine | 38004500 | Critical care (intensivist)      | Provider Specialty | Specialty    
Cytopathology |45756766|Cytopathology |Provider Specialty |  ABMS
Dermatology  |38004452 | Dermatology                        | Provider Specialty | Specialty    
Dermatopathology |45756767|Dermatopathology|Provider Specialty |  ABMS
Developmental-Behavioral Pediatrics |45756768|Developmental-Behavioral Pediatrics|Provider Specialty |  ABMS
Diagnostic Radiology |45756769|Diagnostic Radiology|Provider Specialty |  ABMS
Emergency Medical Services |45756770|Emergency Medical Services|Provider Specialty |  ABMS
Emergency Medicine | 38004510 | Emergency Medicine         | Provider Specialty | Specialty   
Endocrinology, Diabetes and Metabolism |45756771  |Endocrinology, Diabetes and Metabolism|Provider Specialty |  ABMS
Epilepsy |45756772   |Epilepsy|Provider Specialty |  ABMS
General Family Medicine | 38004453 | Family Practice                           | Provider Specialty | Specialty   
Female Pelvic Medicine and Reconstructive Surgery|45756773 |Female Pelvic Medicine and Reconstructive Surgery| Provider Specialty | ABMS
Forensic Psychiatry |45756775|Forensic Psychiatry|Provider Specialty |  ABMS
Gastroenterology |38004455 | Gastroenterology                    | Provider Specialty | Specialty   
General Pediatrics |38004477 | Pediatric Medicine            | Provider Specialty | Specialty    
Geriatric Medicine | 38004478 | Geriatric Medicine                   | Provider Specialty | Specialty   
Geriatric Psychiatry |45756776|Geriatric Psychiatry|Provider Specialty |  ABMS
Gynecologic Oncology |38004513 | Gynecology/Oncology              | Provider Specialty | Specialty    
Hematology | 38004501 | Hematology                           | Provider Specialty | Specialty   
Hospice and Pallative Medicine |45756777 |Hospice and Pallative Medicine|Provider Specialty |  ABMS
Infectious Disease | 38004484 | Infectious Disease                | Provider Specialty | Specialty   
General Internal Medicine | 38004456 | Internal Medicine| Provider Specialty | Specialty   
Internal Medicine - Critical Care Medicine |45756778|Internal Medicine - Critical Care Medicine|Provider Specialty |  ABMS
Interventional Cardiology |45756779 |Interventional Cardiology|Provider Specialty |  ABMS
Interventional Radiology and Diagnostic Radiology |38004511 | Interventional Radiology | Provider Specialty | Specialty   
Maternal and Fetal Medicine |45756780|Maternal and Fetal Medicine |Provider Specialty |  ABMS
Medical Biochemical Genetics |45756781|Medical Biochemical Genetics|Provider Specialty |  ABMS
Medical Genetics and Genomics |45756782 |Medical Genetics and Genomics|Provider Specialty |  ABMS
Medical Oncology |38004507 | Medical Oncology | Provider Specialty | Specialty   
Medical Physics |45756783|Medical Physics|Provider Specialty |  ABMS
Medical Toxicology|45756784|Medical Toxicology|Provider Specialty |  ABMS
Molecular Genetic Pathology |45756785|Molecular Genetic Pathology|Provider Specialty |  ABMS
Neonatal-Perinatal Medicine |45756786|Neonatal-Perinatal Medicine|Provider Specialty |  ABMS
Nephrology | 38004479 | Nephrology                   | Provider Specialty | Specialty   
Neurodevelopmental Disabilities |45756787|Neurodevelopmental Disabilities|Provider Specialty |  ABMS
Neurological Surgery | 38004459 | Neurosurgery                    | Provider Specialty | Specialty   
General Neurology | 38004458 | Neurology                                     | Provider Specialty | Specialty   
Neurology with Special Qualification in Child Neurology |45756788|Neurology with Special Qualification in Child Neurology|Provider Specialty |  ABMS
Neuromuscular Medicine |45756789|Neuromuscular Medicine|Provider Specialty |  ABMS
Neuropathology |45756790 |Neuropathology|Provider Specialty |  ABMS
Neuroradiology |45756791|Neuroradiology|Provider Specialty |  ABMS
Neurotology |45756792|Neurotology|Provider Specialty |  ABMS
Nuclear Medicine |38004476 | Nuclear Medicine                  | Provider Specialty | Specialty   
Nuclear Radiology |45756793 |Nuclear Radiology|Provider Specialty |  ABMS
Obstetrics and Gynecology | 38004461 | Obstetrics/Gynecology              | Provider Specialty | Specialty   
Occupational Medicine |38004492 | Occupational Therapy              | Provider Specialty | Specialty   
Ophthalmology |  38004463 | Ophthalmology                       | Provider Specialty | Specialty     
Orthopaedic Sports Medicine |45756794|Orthopaedic Sports Medicine|Provider Specialty |  ABMS
Orthopaedic Surgery |38004465 | Orthopedic Surgery                | Provider Specialty | Specialty   
Otolaryngology | 38004449 | Otolaryngology                           | Provider Specialty | Specialty   
Pain Medicine | 38004494 | Pain Management                           | Provider Specialty | Specialty   
Pathology |38004466 | Pathology                                  | Provider Specialty | Specialty   
Pathology - Anatomic |45756795|Pathology - Anatomic| Provider Specialty | ABMS
Pathology - Chemical |45756796|Pathology - Chemical | Provider Specialty | ABMS
Pathology - Clinical |45756797|Pathology - Clinical| Provider Specialty | ABMS
Pathology - Forensic |45756798|Pathology - Forensic| Provider Specialty | ABMS
Pathology - Hematology |45756799|Pathology - Hematology| Provider Specialty | ABMS
Pathology - Medical Microbiology |45756800 |Pathology - Medical Microbiology| Provider Specialty | ABMS
Pathology - Molecular Genetic |45756801|Pathology - Molecular Genetic| Provider Specialty | ABMS
Pathology - Pediatric |45756802|Pathology - Pediatric | Provider Specialty | ABMS
Pathology-Anatomic/Pathology-Clinical |45756803|Pathology-Anatomic/Pathology-Clinical | Provider Specialty | ABMS
General Pediatrics | 38004477 | Pediatric Medicine               | Provider Specialty | Specialty   
Pediatric Anesthesiology |45756804|Pediatric Anesthesiology| Provider Specialty | ABMS
Pediatric Cardiology |45756805|Pediatric Cardiology | Provider Specialty | ABMS
Pediatric Critical Care Medicine |45756806|Pediatric Critical Care Medicine| Provider Specialty | ABMS
Pediatric Dermatology |45756807|Pediatric Dermatology| Provider Specialty | ABMS
Pediatric Emergency Medicine |45756808|Pediatric Emergency Medicine| Provider Specialty | ABMS
Pediatric Endocrinology |45756809|Pediatric Endocrinology| Provider Specialty | ABMS
Pediatric Gastroenterology |45756810 |Pediatric Gastroenterology | Provider Specialty | ABMS
Pediatric Hematology-Oncology |45756811 |Pediatric Hematology-Oncology| Provider Specialty | ABMS
Pediatric Infectious Diseases |45756812 |Pediatric Infectious Diseases| Provider Specialty | ABMS
Pediatric Nephrology |45756813 |Pediatric Nephrology| Provider Specialty | ABMS 
Pediatric Otolaryngology |45756814|Pediatric Otolaryngology | Provider Specialty | ABMS
Pediatric Pulmonology |45756815|Pediatric Pulmonology | Provider Specialty | ABMS
Pediatric Radiology |45756816|Pediatric Radiology| Provider Specialty | ABMS
Pediatric Rehabilitation Medicine |45756817|Pediatric Rehabilitation Medicine| Provider Specialty | ABMS
Pediatric Rheumatology |45756818|Pediatric Rheumatology| Provider Specialty | ABMS
Pediatric Surgery |45756819|Pediatric Surgery| Provider Specialty | ABMS
Pediatric Transplant Hepatology |45756820|Pediatric Transplant Hepatology | Provider Specialty | ABMS
Pediatric Urology|45756821|Pediatric Urology| Provider Specialty | ABMS
Physical Medicine and Rehabilitation |38004468 | Physical Medicine And Rehabilitation | Provider Specialty | Specialty   
Plastic Surgery | 38004467 | Plastic And Reconstructive Surgery  | Provider Specialty | Specialty    
Plastic Surgery Within the Head and Neck |45756822|Plastic Surgery Within the Head and Neck| Provider Specialty | ABMS
Preventative Medicine | 38004503 | Preventive Medicine                | Provider Specialty | Specialty   
Psychiatry |38004469 | Psychiatry                             | Provider Specialty | Specialty   
Psychosomatic Medicine |45756823 |Psychosomatic Medicine | Provider Specialty | ABMS
Public Health and General Preventive Medicine |45756824|Public Health and General Preventive Medicine| Provider Specialty | ABMS
Pulmonary Disease | 38004472 | Pulmonary Disease         | Provider Specialty | Specialty   
Radiation Oncology |38004509 | Radiation Oncology    | Provider Specialty | Specialty   
Radiology |45756825|Radiology| Provider Specialty | ABMS
Reproductive Endocrinology/Infertility |45756826|Reproductive Endocrinology/Infertility| Provider Specialty | ABMS
Rheumatology |38004491 | Rheumatology                  | Provider Specialty | Specialty   
Sleep Medicine |45756827|Sleep Medicine| Provider Specialty | ABMS
Spinal Cord Injury Medicine| concept id requested |Spinal Cord Injury Medicine | Provider Specialty | ABMS|
Sports Medicine |45756828|Sports Medicine| Provider Specialty | ABMS
General Surgery | 38004447 | General Surgery            | Provider Specialty | Specialty   
Surgery of the Hand | 38004480 | Hand Surgery                  | Provider Specialty | Specialty   
Surgical Critical Care |45756829|Surgical Critical Care | Provider Specialty | ABMS
Thoracic Surgery | 38004473 | Thoracic Surgery                  | Provider Specialty | Specialty     
Thoracic and Cardiac Surgery |45756830|Thoracic and Cardiac Surgery| Provider Specialty | ABMS
Transplant Hepatology |45756831|Transplant Hepatology| Provider Specialty | ABMS
Undersea and Hyperbaric Medicine |45756832 |Undersea and Hyperbaric Medicine| Provider Specialty | ABMS
Urology | 38004474 | Urology                                   | Provider Specialty | Specialty   
Vascular and Interventional Radiology |45756833|Vascular and Interventional Radiology| Provider Specialty | ABMS
Vascular Neurology |45756834|Vascular Neurology| Provider Specialty | ABMS
Vascular Surgery | 38004496 | Vascular Surgery           | Provider Specialty | Specialty   

* * *
